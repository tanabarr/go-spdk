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
package nvme

import (
	"fmt"
	"testing"

	"../spdk"
)

func checkFailure(shouldSucceed bool, err error) (rErr error) {
	switch {
	case err != nil && shouldSucceed:
		rErr = fmt.Errorf("expected test to succeed, failed unexpectedly: %v", err)
	case err == nil && !shouldSucceed:
		rErr = fmt.Errorf("expected test to fail, succeeded unexpectedly")
	}

	return
}

func TestDiscover(t *testing.T) {
	tests := []struct {
		lib           string
		shouldSucceed bool
	}{
		{
			shouldSucceed: true,
		},
		//{
		//	shouldSucceed: true,
		//},
	}

	for i, tt := range tests {
		if err := spdk.InitSPDKEnv(); err != nil {
			t.Fatal(err.Error())
		}
		//var entries []Namespace
		// err := Update(0, "")
		_, _, err := Discover()
		if checkFailure(tt.shouldSucceed, err) != nil {
			t.Errorf("case %d: %v", i, err)
		}

		// if err := spdk.InitSPDKEnv(); err != nil {
		// t.Fatal(err.Error())
		// }
		//_, _, err = Discover()
		err = Update(0, "")
		if checkFailure(tt.shouldSucceed, err) != nil {
			t.Errorf("case %d: %v", i, err)
		}

		Cleanup()

		//fmt.Println("Discovered NVMe devices: ")
		//for _, e := range entries {
		//	fmt.Printf(
		//		"controller (model/serial): %v/%v, namespace: %v, size: %v\n",
		//		e.CtrlrModel, e.CtrlrSerial, e.Id, e.Size)
		//}

		// if tt.shouldSucceed && len != expLen {
		//			t.Errorf("case %d: expected length %d, got %d", i, expLen, len)
		//		}
	}
}
