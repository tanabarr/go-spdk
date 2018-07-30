#include <stdlib.h>
#include <stdio.h>
#include "nvme_discover.h"

static struct entry_t *g_entries = NULL;

struct entry_t* nvme_discover(void) 
{
	printf("Initializing NVMe Controllers\n");

	struct entry_t *entry_p = NULL; 
	for (int i=0; i < 2 ; i++) {
		entry_p = malloc(sizeof(struct entry_t));
		if (entry_p == NULL) {
			perror("entry malloc");
			exit(1);
		}

		snprintf(
			entry_p->name1,
			sizeof(entry_p->name1),
			"test name one (%d)", 
			i
		);
		snprintf(
			entry_p->name2,
			sizeof(entry_p->name2),
			"test name two (%d)", 
			i
		);
		entry_p->next = g_entries;
		g_entries = entry_p;
	}
//    snprintf(entry_p->name1, sizeof(entry_p->name1), "test name one");
//	snprintf(entry_p->name2, sizeof(entry_p->name2), "test name two");

	printf("Initialization complete.\n");
	//return entry_p;
	return g_entries;
}
