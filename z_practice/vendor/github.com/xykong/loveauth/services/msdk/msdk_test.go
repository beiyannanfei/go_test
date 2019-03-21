package msdk

import (
	"reflect"
	"testing"
	"time"
	"github.com/bouk/monkey"
	"github.com/xykong/loveauth/errors"
)

func TestMsdkKey(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"not null", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MsdkKey(); got == tt.want {
				t.Errorf("MsdkKey() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestHost(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"not null", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Host(); got == tt.want {
				t.Errorf("Host() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestAppId(t *testing.T) {
	type args struct {
		platform Platform
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not null", args{QQ}, ""},
		{"not null", args{Wechat}, ""},
		{"not null", args{Guest}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppId(tt.args.platform); got == tt.want {
				t.Errorf("AppId() = %v, not want %v", got, tt.want)
			}
		})
	}

	if AppId(Guest) != "G_"+AppId(QQ) {
		t.Errorf("AppId Guest platform (%v) should equal G_QQAppId (%v)", AppId(Guest), AppId(QQ))
	}
}

func TestAppKey(t *testing.T) {
	type args struct {
		platform Platform
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not null", args{QQ}, ""},
		{"not null", args{Wechat}, ""},
		{"not null", args{Guest}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppKey(tt.args.platform); got == tt.want {
				t.Errorf("AppKey() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestPlatformId(t *testing.T) {
	type args struct {
		platform Platform
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not null", args{QQ}, ""},
		{"not null", args{Wechat}, ""},
		{"not null", args{Guest}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlatformId(tt.args.platform); got == tt.want {
				t.Errorf("PlatformId() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestSig(t *testing.T) {
	type args struct {
		encode    int
		platform  Platform
		timestamp uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not null", args{2, QQ, 1515256118}, "8c2baab6fe023d415a683e6ddfbadf4d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sig(tt.args.encode, tt.args.platform, tt.args.timestamp); got != tt.want {
				t.Errorf("Sig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlParams(t *testing.T) {

	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	type args struct {
		openId   string
		platform Platform
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not null", args{"61DD81D1162F06DD30924A2183FE58B8", QQ}, "timestamp=138157323&appid=1106545613&sig=ba3a0efab23a67bbc01f48ef21dcbac4&openid=61DD81D1162F06DD30924A2183FE58B8&encode=2"},
		{"not null", args{"oNB9uwqkkIJKlEogST3ehnXLRxGA", Wechat}, "timestamp=138157323&appid=wx6dfb1269f082c9b7&sig=ba3a0efab23a67bbc01f48ef21dcbac4&openid=oNB9uwqkkIJKlEogST3ehnXLRxGA&encode=2"},
		{"not null", args{"G_8f84d850d8792ab39d76d1abf1111528", Guest}, "timestamp=138157323&appid=G_1106545613&sig=ba3a0efab23a67bbc01f48ef21dcbac4&openid=G_8f84d850d8792ab39d76d1abf1111528&encode=2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlParams(tt.args.openId, tt.args.platform); got != tt.want {
				t.Errorf("UrlParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlRequest(t *testing.T) {

	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	type args struct {
		openId   string
		platform Platform
		module   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not null", args{"61DD81D1162F06DD30924A2183FE58B8", QQ, "auth/verify_login"}, "http://msdktest.qq.com/auth/verify_login?timestamp=138157323&appid=1106545613&sig=ba3a0efab23a67bbc01f48ef21dcbac4&openid=61DD81D1162F06DD30924A2183FE58B8&encode=2"},
		{"not null", args{"oNB9uwqkkIJKlEogST3ehnXLRxGA", Wechat, "auth/verify_login"}, "http://msdktest.qq.com/auth/verify_login?timestamp=138157323&appid=wx6dfb1269f082c9b7&sig=ba3a0efab23a67bbc01f48ef21dcbac4&openid=oNB9uwqkkIJKlEogST3ehnXLRxGA&encode=2"},
		{"not null", args{"G_8f84d850d8792ab39d76d1abf1111528", Guest, "auth/verify_login"}, "http://msdktest.qq.com/auth/verify_login?timestamp=138157323&appid=G_1106545613&sig=ba3a0efab23a67bbc01f48ef21dcbac4&openid=G_8f84d850d8792ab39d76d1abf1111528&encode=2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlRequest(tt.args.openId, tt.args.platform, tt.args.module); got != tt.want {
				t.Errorf("UrlRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostRequestRaw(t *testing.T) {
	type args struct {
		openId   string
		platform Platform
		module   string
		body     map[string]interface{}
	}
	var tests []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostRequestRaw(tt.args.openId, tt.args.platform, tt.args.module, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostRequestRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostRequestRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostRequest(t *testing.T) {
	type args struct {
		openId   string
		platform Platform
		module   string
		body     map[string]interface{}
	}
	var tests []struct {
		name string
		args args
		want errors.Code
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := PostRequest(tt.args.openId, tt.args.platform, tt.args.module, tt.args.body); got != tt.want {
				t.Errorf("PostRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
