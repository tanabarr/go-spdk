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
	struct ctrlr_entry	*next;
	char			name[1024];
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
	const struct spdk_nvme_ctrlr_data *cdata = spdk_nvme_ctrlr_get_data(ctrlr);

	entry = malloc(sizeof(struct ctrlr_entry));
	if (entry == NULL) {
		perror("ctrlr_entry malloc");
		exit(1);
	}

	snprintf(entry->name, sizeof(entry->name), "%-20.20s (%-20.20s)", cdata->mn, cdata->sn);

	entry->ctrlr = ctrlr;
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

struct ret_t*
cleanup(bool success)
{
	struct ns_entry *ns_entry = g_namespaces;
	struct ctrlr_entry *ctrlr_entry = g_controllers;
	const struct spdk_nvme_ctrlr_data *cdata;
	struct ns_t *ns;
	struct ctrlr_t *ctrlr;

	while (ns_entry) {
		if (success == true) {
		    ns = malloc(sizeof(struct ns_t));
		    if (ns == NULL) {
		    	perror("ns_t malloc");
		    	exit(1);
		    }

			cdata = spdk_nvme_ctrlr_get_data(ns_entry->ctrlr);

			ns->id = spdk_nvme_ns_get_id(ns_entry->ns);
	        // capacity in GBytes
			ns->size = spdk_nvme_ns_get_size(ns_entry->ns) / 1000000000;
			ns->ctrlr_id = cdata->cntlid;
		    ns->next = g_ns;
		    g_ns = ns;
		}

		struct ns_entry *next = ns_entry->next;
		free(ns_entry);
		ns_entry = next;
	}

	while (ctrlr_entry) {
		if (success == true) {
		    ctrlr = malloc(sizeof(struct ctrlr_t));
		    if (ctrlr == NULL) {
				perror("ctrlr_t malloc");
				exit(1);
		    }
			cdata = spdk_nvme_ctrlr_get_data(ctrlr_entry->ctrlr);
			ctrlr->id = cdata->cntlid;
			snprintf(
				ctrlr->model,
				sizeof(cdata->mn) + 1,
				"%-20.20s",
				cdata->mn
			);
			snprintf(
				ctrlr->serial,
				sizeof(cdata->sn) + 1,
				"%-20.20s",
				cdata->sn
			);
			snprintf(
				ctrlr->pci_addr,
				sizeof(ctrlr->pci_addr) + 1,
				"%04x:%02x:%02x.%02x",
		        dev->pci_addr.domain, dev->pci_addr.bus,
				dev->pci_addr.dev, dev->pci_addr.func;
			);
		    ctrlr->next = g_ctrlr;
		    g_ctrlr = ctrlr;
		}

		struct ctrlr_entry *next = ctrlr_entry->next;

		spdk_nvme_detach(ctrlr_entry->ctrlr);
		free(ctrlr_entry);
		ctrlr_entry = next;
	}

	ret = malloc(sizeof(struct ret_t));
	ret->success = success;
	ret->nss = g_ns;
	ret->ctrlrs = g_ctrlr;

	return ret;
}

struct ret_t* nvme_discover(void)
{
	int rc;

	/*
	 * Start the SPDK NVMe enumeration process.  probe_cb will be called
	 *  for each NVMe controller found, giving our application a choice on
	 *  whether to attach to each controller.  attach_cb will then be
	 *  called for each controller after the SPDK NVMe driver has completed
	 *  initializing the controller we chose to attach.
	 */
	rc = spdk_nvme_probe(NULL, NULL, probe_cb, attach_cb, NULL);
	if (rc != 0) {
		fprintf(stderr, "spdk_nvme_probe() failed\n");
		return cleanup(false);
	}

	if (g_controllers == NULL) {
		fprintf(stderr, "no NVMe controllers found\n");
		return cleanup(false);
	}

	return cleanup(true);
}
