#ifndef NVMEDISCOVER_H
#define NVMEDISCOVER_H
struct ns_t {
    int        id;
    char        ctrlr_model[1024];
    char        ctrlr_serial[1024];
    int         size;
    struct ns_t *next;
};
struct ns_t* nvme_discover(void);
#endif
