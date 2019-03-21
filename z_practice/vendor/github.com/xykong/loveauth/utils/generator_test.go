package utils

import (
	"testing"

	"github.com/xykong/snowflake"
	"github.com/sirupsen/logrus"
)

func Test_initSnowflake(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}
	for range tests {
		t.Run("", func(t *testing.T) {
			initSnowflake()
		})
	}
}

func TestGenerateId(t *testing.T) {
	tests := []struct {
		name string
		want snowflake.ID
	}{
		{"", 160469069730},
		{"", 160469069730},
		{"", 160469069730},
		{"", 160469069730},
		{"", 160469069730},
		{"", 160469069730},
		{"", 160469069730},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateId()
			t.Logf("GenerateId: %v", got)
			if got < tt.want {
				t.Errorf("GenerateId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdString(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{0}, "A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IdString(tt.args.id); got != tt.want {
				t.Errorf("IdString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateAuthCode(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateAuthCode(); len(got) != tt.want {
				t.Errorf("GenerateAuthCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateAuthCodeRaw(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"", args{1}, 6},
		{"", args{12}, 6},
		{"", args{123}, 6},
		{"", args{1234}, 6},
		{"", args{12345}, 6},
		{"", args{123456}, 6},
		{"", args{1234567}, 6},
		{"", args{12345678}, 6},
		{"", args{123456789}, 6},
		{"", args{1234567890}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateAuthCodeRaw(tt.args.code); len(got) != tt.want {
				t.Errorf("generateAuthCodeRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePaymentSequence(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePaymentSequence(1)

			logrus.Infof("GeneratePaymentSequence: %v", got)

			if got == tt.want {
				t.Errorf("GeneratePaymentSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}
