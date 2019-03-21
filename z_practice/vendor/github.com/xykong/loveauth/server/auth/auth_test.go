package auth

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage/model"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		globalId          int64
		request           model.DoAuthRequest
		expirationSeconds int64
		salt              string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(tt.args.globalId, tt.args.request, tt.args.expirationSeconds, tt.args.salt); got != tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUse(t *testing.T) {
	type args struct {
		aps []Provider
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Use(tt.args.aps...)
		})
	}
}

func TestStart(t *testing.T) {
	type args struct {
		group *gin.RouterGroup
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start(tt.args.group)
		})
	}
}

func TestAuthenticator_Auth(t *testing.T) {
	type fields struct {
		Provider Provider
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Authenticator{
				Provider: tt.fields.Provider,
			}
			o.Auth(tt.args.c)
		})
	}
}

func Test_checkAccountBanState(t *testing.T) {
	type args struct {
		account *model.Account
	}
	tests := []struct {
		name  string
		args  args
		want  model.AccountState
		want1 string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkAccountBanState(tt.args.account)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkAccountBanState() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkAccountBanState() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_isWhiteListContains(t *testing.T) {
	type args struct {
		request *model.DoAuthRequest
		account *model.Account
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "219.141.227.166",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "1.1.1.1",},}, nil}, false},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "8.8.8.8",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "8.8.8.1",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "8.8.8.18",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "127.0.0.1",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "10.0.0.1",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "::1",},}, nil}, true},
		{"", args{&model.DoAuthRequest{Extra: model.Extra{VClientIP: "2001:db8::1",},}, nil}, true},
	}

	settings.Set("loveauth_white_list", "ip", []string{
		"219.141.227.166",
		"8.8.8.8/24",
		"10.0.0.1/24",
		"127.0.0.1",
		"2001:db8::/32",
		"::1",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWhiteListContains(tt.args.request); got != tt.want {
				t.Errorf("isWhiteListContains(%v) = %v, want %v", tt.args.request, got, tt.want)
			}
		})
	}
}
