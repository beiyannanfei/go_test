package midas

import (
	"github.com/xykong/loveauth/errors"
	"reflect"
	"testing"
)

func TestCancelPay(t *testing.T) {
	type args struct {
		request DoMidasCancelPayReq
	}

	//noinspection SpellCheckingInspection
	tests := []struct {
		name  string
		args  args
		want  *DoMidasCancelPayRsp
		want1 error
	}{
		//{"ios_qq", args{request: DoMidasCancelPayReq{
		//	OpenId:  "C64B7ECD63B5BAF7E1AAC4984777E89F",
		//	OpenKey: "7A7CD73BE0747E45A8617EADFD7770E9",
		//	Pf:      "qq_m_qq-73213123-iap-1001-qq-1106545613-C64B7ECD63B5BAF7E1AAC4984777E89F",
		//	Pfkey:   "9f0e39f84ca2eddb43689dc3600266fa",
		//	UserIp:  "10.1.16.69",
		//	Amount:  1,
		//	BillNo:  "1234567890",
		//	PayItem: "test-item",
		//}}, nil, errors.NewCodeString(errors.AuthTokenInvalid, "token校验失败(18)")},
		//{"ios_wechat", args{request: DoMidasCancelPayReq{
		//	OpenId:  "oNB9uwjqf-ufdJjUI5RkJCHSlcPU",
		//	OpenKey: "13_wdlf1OZhO0432WxWO_CVQKgb3CPh8UHE8fBbjJBDJu1214Y3S3eES4F3Ky1IpY2oCbeSdiG9rx9QDFiP2GIKhczxsBUo9a9EdjVRValwFNI",
		//	Pf:      "wechat_wx-1001-iap-1001-wx-wx6dfb1269f082c9b7-oNB9uwjqf-ufdJjUI5RkJCHSlcPU",
		//	Pfkey:   "c6567a45b711dca73ab6011a7f361fee",
		//	UserIp:  "10.1.16.69",
		//	Amount:  1,
		//	BillNo:  "1234567890",
		//	PayItem: "test-item",
		//}}, nil, errors.NewCodeString(errors.AuthTokenInvalid, "token校验失败(18)")},
		//{"ios_qq_test", args{request: DoMidasCancelPayReq{
		//	OpenId:  "9FCD2AFE94552788E8F531111B664B4D",
		//	OpenKey: "E5FB131A48F710A27683991CDB3EA71E",
		//	Pf:      "qq_m_qq-1001-iap-1001-qq-1106545613-9FCD2AFE94552788E8F531111B664B4D",
		//	Pfkey:   "347d13448bb4ec85134b0672c06c436a",
		//	UserIp:  "10.1.16.69",
		//	Amount:  1,
		//	BillNo:  "1234567890",
		//	PayItem: "test-item",
		//}}, nil, errors.NewCodeString(errors.AuthTokenInvalid, "token校验失败(18)")},
		//{"android_qq", args{request: DoMidasCancelPayReq{
		//	OpenId:  "C64B7ECD63B5BAF7E1AAC4984777E89F",
		//	OpenKey: "696BE3AEF2DCEABE2823484BEB6240C0",
		//	Pf:      "desktop_m_qq-73213123-android-73213123-qq-1106545613-C64B7ECD63B5BAF7E1AAC4984777E89F",
		//	Pfkey:   "73f775e3d19dce894ca2524cf7922a6a",
		//	UserIp:  "10.1.16.69",
		//	Amount:  1,
		//	BillNo:  "1234567890",
		//	PayItem: "test-item",
		//	PlatId:  1,
		//}}, nil, errors.NewCodeString(errors.AuthTokenInvalid, "token校验失败(18)")},
		{"android_qq", args{request: DoMidasCancelPayReq{
			OpenId:  "oNB9uwjqf-ufdJjUI5RkJCHSlcPU",
			OpenKey: "13_G_WMcf_vCh_a5mxL0M_kYekC3jhvCynjnWKKldQM5VkzHcO7QasXG56to1YwwKmq3yPzavAj5SHHUmi1TxEq4o7EShgEPeG82VOGP8b2fVE",
			Pf:      "desktop_m_wx-1001-android-73213123-wx-wx6dfb1269f082c9b7-oNB9uwjqf-ufdJjUI5RkJCHSlcPU",
			Pfkey:   "5ab50183e7cc2a95970e211e1f986f8a",
			UserIp:  "10.1.16.69",
			Amount:  1,
			BillNo:  "1234567890",
			PayItem: "test-item",
			PlatId:  1,
		}}, nil, errors.NewCodeString(errors.AuthTokenInvalid, "token校验失败(18)")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CancelPay(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CancelPay() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CancelPay() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
