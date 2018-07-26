#include <stdlib.h>
#include <stdio.h>
#include "nvme_discover.h"

//static struct entry_t *g_entries = NULL;

//struct entr nvme_discover(void) 
struct entry_t* nvme_discover(void) 
{
	printf("Initializing NVMe Controllers\n");

	struct entry_t *entry_p;

	entry_p = malloc(sizeof(struct entry_t));
	if (entry_p == NULL) {
		perror("entry malloc");
		exit(1);
	}

	entry_p->name1 = "test name one";
	entry_p->name2 = "test name two";
//    snprintf(entry_p->name1, sizeof(entry_p->name1), "test name one");
//	snprintf(entry_p->name2, sizeof(entry_p->name2), "test name two");

	printf("Initialization complete.\n");
	//return entry_p;
	return entry_p;
}
