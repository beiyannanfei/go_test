package midas

import (
	"reflect"
	"testing"
	"github.com/xykong/loveauth/errors"
)

func TestBuyGoods(t *testing.T) {
	type args struct {
		request DoMidasBuyGoodsReq
	}

	//noinspection SpellCheckingInspection
	tests := []struct {
		name  string
		args  args
		want  *DoMidasBuyGoodsRsp
		want1 error
	}{
		{"android", args{request: DoMidasBuyGoodsReq{
			OpenId:  "3DD5817D1EACA32E94B1E838823E0A5E",
			OpenKey: "603A4FF276AAEA100AA519DDE79A8074",
			Pf:      "desktop_m_qq-73213123-android-73213123-qq-1106545613-3DD5817D1EACA32E94B1E838823E0A5E",
			Pfkey:   "04376e20bc860f9c5c6940ae6d0876ae",
			UserIp:  "10.1.16.111",
			PlatId:  1,

			PayItem:     "10001*10*1",
			GoodsMeta:   "物品10001*物品10001介绍+描述",
			AppMode:     1,
			Amount:      10,
			AppMetadata: "sendfromloveauth",
		}}, nil, errors.NewCodeString(errors.AuthTokenInvalid, "token校验失败(18)")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := BuyGoods(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuyGoods() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("BuyGoods() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
