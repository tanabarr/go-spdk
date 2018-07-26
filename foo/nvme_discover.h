#ifndef NVMEDISCOVER_H
#define NVMEDISCOVER_H
struct entry_t {
	char			*name1;
	char			*name2;
};
struct entry_t* nvme_discover(void);
#endif
