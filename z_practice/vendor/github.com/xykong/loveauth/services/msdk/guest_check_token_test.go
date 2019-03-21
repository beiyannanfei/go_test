package msdk

import (
	"reflect"
	"testing"
)

//2018/06/15 18:58:27 INFO[2018-06-15T18:58:27+08:00] Msdk post request.
//	body="{\"accessToken\":\"XiYPfi9kQkfRwgDL6zVkeSqNKebn1D6elZBPw9LldfU=\",\"guestid\":\"\"}"
//	elapsed=15.130878ms
//	request="http://msdktest.qq.com//auth/guest_check_token?timestamp=1529060307&appid=1106545613&sig=c789769876b0eb6e3e163dfd3fd94ad4&openid=G_111898318f21364b4ab757a15cd9d038&encode=2"
//	response="{\"err_type\":\"5\",\"msg\":\"input param is null.\",\"ret\":-1011}\n"

func TestGuestCheckToken(t *testing.T) {
	type args struct {
		openId string
		token  string
	}
	tests := []struct {
		name    string
		args    args
		want    *DoGuestCheckTokenRsp
		wantErr bool
	}{
	//{"", args{"G_111898318f21364b4ab757a15cd9d038", "XiYPfi9kQkfRwgDL6zVkeSqNKebn1D6elZBPw9LldfU="}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GuestCheckToken(tt.args.openId, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuestCheckToken(%v, %v) error = %v, wantErr %v", tt.args.openId, tt.args.token, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GuestCheckToken(%v, %v) = %v, want %v", tt.args.openId, tt.args.token, got, tt.want)
			}
		})
	}
}
