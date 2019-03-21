package midas

import (
	"reflect"
	"testing"
)

//noinspection ALL
func Test_makeSig(t *testing.T) {
	type args struct {
		method  string
		urlPath string
		params  map[string]interface{}
		secret  string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{method: "GET", urlPath: "/v3/user/get_info", params: map[string]interface{}{
			"openid":  "11111111111111111",
			"openkey": "2222222222222222",
			"appid":   "123456",
			"pf":      "qzone",
			"format":  "json",
			"userip":  "112.90.139.30",
		}, secret: "228bf094169a40a3bd188ba37ebe8723&"}, "FdJkiDYwMj5Aj1UG2RUPc83iokk="},
		{"", args{method: "GET", urlPath: "/mpay/pay_m", params: map[string]interface{}{
			"accounttype": "common",
			"amt":         "1",
			"appid":       "1450014200",
			"appremark":   "midas pay remark",
			"billno":      "1234567890",
			"format":      "json",
			"openid":      "6642C4E6730B334EE5D8544A12D41825",
			"openkey":     "CB40CD8A40650F0D1FD3FFED38C49880",
			"payitem":     "test-item",
			"pf":          "desktop_m_qq-73213123-android-73213123-qq-1106545613-6642C4E6730B334EE5D8544A12D41825",
			"pfkey":       "2aec5434911e5112c0cee5ff24bc9a82",
			"ts":          "1522292397",
			"userip":      "10.1.16.111",
			"zoneid":      "1",
		}, secret: "n32dMJcQYsWsjiVvuZx7AnEI0MQf6HVR&"}, "OQGafxBXm3p2+/n16LmKyJOMmMY="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeSig(tt.args.method, tt.args.urlPath, tt.args.params, tt.args.secret); got != tt.want {
				t.Errorf("makeSig() = %v, want %v", got, tt.want)
			}
		})
	}
}

//noinspection ALL
func Test_makeSource(t *testing.T) {
	type args struct {
		method  string
		urlPath string
		params  map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"",
			args{
				method: "GET", urlPath: "/v3/user/get_info",
				params: map[string]interface{}{
					"openid":  "11111111111111111",
					"openkey": "2222222222222222",
					"appid":   "123456",
					"pf":      "qzone",
					"format":  "json",
					"userip":  "112.90.139.30",
				},
			},
			"GET&%2Fv3%2Fuser%2Fget_info&appid%3D123456%26format%3Djson%26openid%3D11111111111111111%26openkey%3D2222222222222222%26pf%3Dqzone%26userip%3D112.90.139.30",
		},
		{"",
			args{
				method: "GET", urlPath: "/mpay/pay_m",
				params: map[string]interface{}{
					"accounttype": "common",
					"amt":         "1",
					"appid":       "1450014200",
					"appremark":   "midas pay remark",
					"billno":      "1234567890",
					"format":      "json",
					"openid":      "6642C4E6730B334EE5D8544A12D41825",
					"openkey":     "CB40CD8A40650F0D1FD3FFED38C49880",
					"payitem":     "test-item",
					"pf":          "desktop_m_qq-73213123-android-73213123-qq-1106545613-6642C4E6730B334EE5D8544A12D41825",
					"pfkey":       "2aec5434911e5112c0cee5ff24bc9a82",
					"ts":          "1522292397",
					"userip":      "10.1.16.111",
					"zoneid":      "1",
				},
			},
			"GET&%2Fmpay%2Fpay_m&accounttype%3Dcommon%26amt%3D1%26appid%3D1450014200%26appremark%3Dmidas%20pay%20remark%26billno%3D1234567890%26format%3Djson%26openid%3D6642C4E6730B334EE5D8544A12D41825%26openkey%3DCB40CD8A40650F0D1FD3FFED38C49880%26payitem%3Dtest-item%26pf%3Ddesktop_m_qq-73213123-android-73213123-qq-1106545613-6642C4E6730B334EE5D8544A12D41825%26pfkey%3D2aec5434911e5112c0cee5ff24bc9a82%26ts%3D1522292397%26userip%3D10.1.16.111%26zoneid%3D1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeSource(tt.args.method, tt.args.urlPath, tt.args.params); got != tt.want {
				t.Errorf("makeSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

//noinspection ALL
func Test_sortByKey(t *testing.T) {
	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"", args{params: map[string]interface{}{
			"openid":  "11111111111111111",
			"openkey": "2222222222222222",
			"appid":   "123456",
			"pf":      "qzone",
			"format":  "json",
			"userip":  "112.90.139.30",
		}}, []string{"appid", "format", "openid", "openkey", "pf", "userip"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortByKey(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

//noinspection ALL
func Test_urlEncode(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{"", args{"/v3/user/get_info"}, "%2Fv3%2Fuser%2Fget_info"},
		{"", args{" "}, "%20"},
		{"", args{"+"}, "%2B"},
		{"", args{"*"}, "%2A"},
		{"", args{"="}, "%3D"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := urlEncode(tt.args.s); gotResult != tt.wantResult {
				t.Errorf("urlEncode(%v) = %v, want %v", tt.args.s, gotResult, tt.wantResult)
			}
		})
	}
}
