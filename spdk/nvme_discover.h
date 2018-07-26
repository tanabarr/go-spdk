#ifndef NVMEDISCOVER_H
#define NVMEDISCOVER_H
struct ns_t {
    char        *id;
    char        *ctrlr_name;
    int         size;
    struct ns_t *next;
};
struct ns_t* nvme_discover(void);
#endif
