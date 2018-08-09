package spdk

import (
	"fmt"
	"testing"
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

func TestNVMeDiscover(t *testing.T) {
	tests := []struct {
		lib           string
		shouldSucceed bool
	}{
		{
			shouldSucceed: false,
		},
		//{
		//	shouldSucceed: true,
		//},
	}

	for i, tt := range tests {
		//var entries []Namespace
		_, err := NVMeDiscover()
		if checkFailure(tt.shouldSucceed, err) != nil {
			t.Errorf("case %d: %v", i, err)
		}

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
