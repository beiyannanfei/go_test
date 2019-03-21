package utils

import (
	"net"
	"reflect"
	"testing"
)

func Test_get_external(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get_external(); IsPublicIPString(got) != tt.want {
				t.Errorf("get_external() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntranetIP(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetIntranetIP()
		})
	}
}

func TestTabaoAPI(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want *IPInfo
	}{
		{"", args{"8.8.8.8"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TabaoAPI(tt.args.ip); got == nil {
				t.Errorf("TabaoAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inet_ntoa(t *testing.T) {
	type args struct {
		ipnr int64
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inet_ntoa(tt.args.ipnr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inet_ntoa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inet_aton(t *testing.T) {
	type args struct {
		ipnr net.IP
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inet_aton(tt.args.ipnr); got != tt.want {
				t.Errorf("inet_aton() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPBetween(t *testing.T) {
	type args struct {
		test net.IP
		from net.IP
		to   net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPBetween(tt.args.test, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("IPBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPublicIP(t *testing.T) {
	type args struct {
		IP net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPublicIP(tt.args.IP); got != tt.want {
				t.Errorf("IsPublicIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublicIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPublicIP(); net.ParseIP(got) == nil {
				t.Errorf("GetPublicIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPublicIPString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"1.1.1.1"}, true},
		{"", args{"8.8.8.8"}, true},
		{"", args{"0"}, false},
		{"", args{"127.0.0.1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPublicIPString(tt.args.s); got != tt.want {
				t.Errorf("IsPublicIPString() = %v, want %v", got, tt.want)
			}
		})
	}
}
