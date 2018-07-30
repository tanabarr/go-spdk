#ifndef NVMEDISCOVER_H
#define NVMEDISCOVER_H
struct entry_t {
	char			name1[1024];
	char			name2[1024];
	struct entry_t  *next;
};
struct entry_t* nvme_discover(void);
#endif
