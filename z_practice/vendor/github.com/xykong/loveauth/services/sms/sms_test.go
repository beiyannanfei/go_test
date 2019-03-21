package sms

import (
	"testing"
	"github.com/xykong/loveauth/utils"
)

func Test_sendAuth(t *testing.T) {
	type args struct {
		code       string
		mobile     string
		templateId int
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{utils.GenerateAuthCode(), "18810776836", 237117}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendAuth(tt.args.code, tt.args.mobile, tt.args.templateId)
		})
	}
}
