//
// (C) Copyright 2018 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

#include "spdk/stdinc.h"

#include "spdk/nvme.h"
#include "spdk/env.h"

#include "nvme_discover.h"

struct ctrlr_entry {
	struct spdk_nvme_ctrlr	*ctrlr;
	const char *tr_addr;
	struct ctrlr_entry    	*next;
};

struct ns_entry {
	struct spdk_nvme_ctrlr	*ctrlr;
	struct spdk_nvme_ns	*ns;
	struct ns_entry		*next;
	struct spdk_nvme_qpair	*qpair;
};

static struct ctrlr_entry *g_controllers = NULL;
static struct ns_entry *g_namespaces = NULL;
static struct ctrlr_t *g_ctrlr = NULL;
static struct ns_t *g_ns = NULL;

static void
register_ns(struct spdk_nvme_ctrlr *ctrlr, struct spdk_nvme_ns *ns)
{
	struct ns_entry *entry;
	const struct spdk_nvme_ctrlr_data *cdata;

	/*
	 * spdk_nvme_ctrlr is the logical abstraction in SPDK for an NVMe
	 *  controller.  During initialization, the IDENTIFY data for the
	 *  controller is read using an NVMe admin command, and that data
	 *  can be retrieved using spdk_nvme_ctrlr_get_data() to get
	 *  detailed information on the controller.  Refer to the NVMe
	 *  specification for more details on IDENTIFY for NVMe controllers.
	 */
	cdata = spdk_nvme_ctrlr_get_data(ctrlr);

	if (!spdk_nvme_ns_is_active(ns)) {
		printf("Controller %-20.20s (%-20.20s): Skipping inactive NS %u\n",
		       cdata->mn, cdata->sn,
		       spdk_nvme_ns_get_id(ns));
		return;
	}

	entry = malloc(sizeof(struct ns_entry));
	if (entry == NULL) {
		perror("ns_entry malloc");
		exit(1);
	}

	entry->ctrlr = ctrlr;
	entry->ns = ns;
	entry->next = g_namespaces;
	g_namespaces = entry;
}

static bool
probe_cb(void *cb_ctx, const struct spdk_nvme_transport_id *trid,
	 struct spdk_nvme_ctrlr_opts *opts)
{
	return true;
}

static void
attach_cb(void *cb_ctx, const struct spdk_nvme_transport_id *trid,
	  struct spdk_nvme_ctrlr *ctrlr, const struct spdk_nvme_ctrlr_opts *opts)
{
	int nsid, num_ns;
	struct ctrlr_entry *entry;
	struct spdk_nvme_ns *ns;

	entry = malloc(sizeof(struct ctrlr_entry));
	if (entry == NULL) {
		perror("ctrlr_entry malloc");
		exit(1);
	}

	entry->ctrlr = ctrlr;
	entry->tr_addr = trid->traddr;
	entry->next = g_controllers;
	g_controllers = entry;

	/*
	 * Each controller has one or more namespaces.  An NVMe namespace is basically
	 *  equivalent to a SCSI LUN.  The controller's IDENTIFY data tells us how
	 *  many namespaces exist on the controller.  For Intel(R) P3X00 controllers,
	 *  it will just be one namespace.
	 *
	 * Note that in NVMe, namespace IDs start at 1, not 0.
	 */
	num_ns = spdk_nvme_ctrlr_get_num_ns(ctrlr);
	for (nsid = 1; nsid <= num_ns; nsid++) {
		ns = spdk_nvme_ctrlr_get_ns(ctrlr, nsid);
		if (ns == NULL) {
			continue;
		}
		register_ns(ctrlr, ns);
	}
}

static struct ret_t*
cleanup(bool success)
{
	struct ret_t *ret;

	ret = malloc(sizeof(struct ret_t));
	ret->success = success;
	ret->nss = NULL;
	ret->ctrlrs = NULL;

	return ret;
}

static void
cleanup2(void)
{
	struct ns_entry *ns_entry = g_namespaces;
	struct ctrlr_entry *ctrlr_entry = g_controllers;

	printf("inside cleaning\n");

	while (ns_entry) {
		printf("inside cleaning:  nclears, next %s\n", ns_entry, ns_entry->next);
		struct ns_entry *next = ns_entry->next;
		printf("inside cleaning: freeing\n");
		free(ns_entry);
		printf("inside cleaning: assigning next\n");
		ns_entry = next;
	}

	printf("inside cleaning: finished entry starting controllers\n");
	
	while (ctrlr_entry) {
		printf("inside cleaning: controller %s, next %s\n", ctrlr_entry, ctrlr_entry->next);
		struct ctrlr_entry *next = ctrlr_entry->next;
		printf("inside cleaning: detaching controller\n");
		spdk_nvme_detach(ctrlr_entry->ctrlr);
		printf("inside cleaning: freeing controller\n");
		free(ctrlr_entry);
		printf("inside cleaning: assigning next\n");
		ctrlr_entry = next;
	}
	printf("inside cleaning: finished entry starting controllers\n");
}

struct ret_t* nvme_discover(void)
{
	return cleanup(true);
}

int nvme_fwupdate2(int ctrlr_id, char *path)
{
	int rc;
	g_controllers = NULL;
	g_namespaces = NULL;
	
	/*
	 * Start the SPDK NVMe enumeration process.  probe_cb will be called
	 *  for each NVMe controller found, giving our application a choice on
	 *  whether to attach to each controller.  attach_cb will then be
	 *  called for each controller after the SPDK NVMe driver has completed
	 *  initializing the controller we chose to attach.
	 */
	printf("probing\n");
	rc = spdk_nvme_probe(NULL, NULL, probe_cb, attach_cb, NULL);
	if (rc != 0) {
		fprintf(stderr, "spdk_nvme_probe() failed\n");
		cleanup2();
		return 1;
	}
	printf("checking\n");
	if (g_controllers == NULL) {
		fprintf(stderr, "no NVMe controllers found\n");
		cleanup2();
		return 1;
	}

//	printf("looking for controller %d", ctrlr_id);
//	struct ctrlr_entry *ctrlr_entry = g_controllers;
//	const struct spdk_nvme_ctrlr_data *cdata;
//	while (ctrlr_entry) {
//		struct ctrlr_entry *next = ctrlr_entry->next;
//		cdata = spdk_nvme_ctrlr_get_data(ctrlr_entry->ctrlr);
//		printf("found controller %d", cdata->cntlid);
//		free(ctrlr_entry);
//		ctrlr_entry = next;
//	}
	printf("cleaning\n");
	cleanup2();

	g_controllers = NULL;
	g_namespaces = NULL;
	
	printf("probing\n");
	rc = spdk_nvme_probe(NULL, NULL, probe_cb, attach_cb, NULL);
	if (rc != 0) {
		fprintf(stderr, "spdk_nvme_probe() failed\n");
		cleanup2();
		return 1;
	}
	printf("checking\n");
	if (g_controllers == NULL) {
		fprintf(stderr, "no NVMe controllers found\n");
		cleanup2();
		return 1;
	}
	printf("cleaning\n");
	cleanup2();

	printf("exiting\n");
	return rc;
}