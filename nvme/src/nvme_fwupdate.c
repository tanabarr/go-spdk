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

#include "nvme_fwupdate.h"

struct ctrlr_entry {
	struct spdk_nvme_ctrlr	*ctrlr;
	const char *tr_addr;
	struct ctrlr_entry    	*next;
};

static struct ctrlr_entry *g_controllers = NULL;

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
	struct ctrlr_entry *entry;

	entry = malloc(sizeof(struct ctrlr_entry));
	if (entry == NULL) {
		perror("ctrlr_entry malloc");
		exit(1);
	}

	entry->ctrlr = ctrlr;
	entry->tr_addr = trid->traddr;
	entry->next = g_controllers;
	g_controllers = entry;
}

static void
cleanup(void)
{
	struct ctrlr_entry *ctrlr_entry = g_controllers;

	while (ctrlr_entry) {
		struct ctrlr_entry *next = ctrlr_entry->next;

		spdk_nvme_detach(ctrlr_entry->ctrlr);
		free(ctrlr_entry);
		ctrlr_entry = next;
	}

	return;
}

int nvme_fwupdate(int ctrlr_id, char *path)
{
	/*
	 * Start the SPDK NVMe enumeration process.  probe_cb will be called
	 *  for each NVMe controller found, giving our application a choice on
	 *  whether to attach to each controller.  attach_cb will then be
	 *  called for each controller after the SPDK NVMe driver has completed
	 *  initializing the controller we chose to attach.
	 */
	int rc = spdk_nvme_probe(NULL, NULL, probe_cb, attach_cb, NULL);
	if (rc != 0) {
		fprintf(stderr, "spdk_nvme_probe() failed\n");
		cleanup();
		return 1;
	}

	if (g_controllers == NULL) {
		fprintf(stderr, "no NVMe controllers found\n");
		cleanup();
		return 1;
	}

	printf("looking for controller %d", ctrlr_id);
	struct ctrlr_entry *ctrlr_entry = g_controllers;
	const struct spdk_nvme_ctrlr_data *cdata;

	while (ctrlr_entry) {
		struct ctrlr_entry *next = ctrlr_entry->next;

		cdata = spdk_nvme_ctrlr_get_data(ctrlr_entry->ctrlr);
		printf("found controller %d", cdata->cntlid);

		free(ctrlr_entry);
		ctrlr_entry = next;
	}

	cleanup();
	return 0;
}
//static void
//update_firmware_image(void)
//{
//	int					rc;
//	int					fd = -1;
//	int					slot;
//	unsigned int				size;
//	struct stat				fw_stat;
//	char					path[256];
//	void					*fw_image;
//	struct dev				*ctrlr;
//	const struct spdk_nvme_ctrlr_data	*cdata;
//	enum spdk_nvme_fw_commit_action		commit_action;
//	struct spdk_nvme_status			status;
//
//	ctrlr = get_controller();
//	if (ctrlr == NULL) {
//		printf("Invalid controller PCI BDF.\n");
//		return;
//	}
//
//	cdata = ctrlr->cdata;
//
//	if (!cdata->oacs.firmware) {
//		printf("Controller does not support firmware download and commit command\n");
//		return;
//	}
//
//	printf("Please Input The Path Of Firmware Image\n");
//
////	if (get_line(path, sizeof(path), stdin) == NULL) {
////		printf("Invalid path setting\n");
////		while (getchar() != '\n');
////		return;
////	}
////
//	fd = open(path, O_RDONLY);
//	if (fd < 0) {
//		perror("Open file failed");
//		return;
//	}
//	rc = fstat(fd, &fw_stat);
//	if (rc < 0) {
//		printf("Fstat failed\n");
//		close(fd);
//		return;
//	}
//
//	if (fw_stat.st_size % 4) {
//		printf("Firmware image size is not multiple of 4\n");
//		close(fd);
//		return;
//	}
//
//	size = fw_stat.st_size;
//
//	fw_image = spdk_dma_zmalloc(size, 4096, NULL);
//	if (fw_image == NULL) {
//		printf("Allocation error\n");
//		close(fd);
//		return;
//	}
//
//	if (read(fd, fw_image, size) != ((ssize_t)(size))) {
//		printf("Read firmware image failed\n");
//		close(fd);
//		spdk_dma_free(fw_image);
//		return;
//	}
//	close(fd);
//
//	printf("Please Input Slot(0 - 7):\n");
//	if (!scanf("%d", &slot)) {
//		printf("Invalid Slot\n");
//		spdk_dma_free(fw_image);
//		while (getchar() != '\n');
//		return;
//	}
//
//	commit_action = SPDK_NVME_FW_COMMIT_REPLACE_AND_ENABLE_IMG;
//	rc = spdk_nvme_ctrlr_update_firmware(ctrlr->ctrlr, fw_image, size, slot, commit_action, &status);
//	if (rc == -ENXIO && status.sct == SPDK_NVME_SCT_COMMAND_SPECIFIC &&
//	    status.sc == SPDK_NVME_SC_FIRMWARE_REQ_CONVENTIONAL_RESET) {
//		printf("conventional reset is needed to enable firmware !\n");
//	} else if (rc) {
//		printf("spdk_nvme_ctrlr_update_firmware failed\n");
//	} else {
//		printf("spdk_nvme_ctrlr_update_firmware success\n");
//	}
//	spdk_dma_free(fw_image);
//}

//	spdk_env_opts_init(&opts);
//	opts.name = "nvme_manage";
//	opts.core_mask = "0x1";
//	opts.shm_id = g_shm_id;
//	if (spdk_env_init(&opts) < 0) {
//		fprintf(stderr, "Unable to initialize SPDK env\n");
//		return 1;
//	}
