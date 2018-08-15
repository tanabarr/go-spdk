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

#ifndef NVMEDISCOVER_H
#define NVMEDISCOVER_H
#include <stdbool.h>
struct ctrlr_t {
    int             id;
    char            model[1024];
    char            serial[1024];
    char            pci_addr[1024];
    struct ctrlr_t  *next;
};
struct ns_t {
    int             id;
    int             size;
    int             ctrlr_id;
    struct ns_t     *next;
};
struct ret_t {
    bool            success;
    struct ctrlr_t  *ctrlrs;
    struct ns_t     *nss;
};
struct ret_t* nvme_discover(void);
#endif
