package tss_sdk

import (
	"github.com/sirupsen/logrus"
)

/*
#include "tss_sdk_anti.h"

extern void goCallback(void*);

*/
import "C"

//export TssSdkSendDataToClientV3Callback
func TssSdkSendDataToClientV3Callback(antiData *C.TssSdkAntiSendDataInfoV3) C.TssSdkProcResult {

	logrus.Warnf("TssSdkSendDataToClientV3Callback go callback: %+v", *antiData)

	return C.TSS_SDK_PROC_OK
}
