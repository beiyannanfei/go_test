// +build linux

package tss_sdk

import (
	"os"
	"testing"
)

var tssSdk *TssSdk

func init() {
	os.Chdir("../../")
}

//func TestStart(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		{},
//	}
//
//	logrus.Warn(os.Getwd())
//
//	//settings.Set("tencent", "tss.shared_lib_dir", "../../")
//	//settings.Set("tencent", "tss.configs", "../../configs")
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			Start()
//		})
//	}
//
//	time.Sleep(time.Second * 2)
//}

func TestTssSdk_Load(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{},
	}

	tssSdk = &TssSdk{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tssSdk.Load()
		})
	}
}

func TestTssSdk_AddUser(t *testing.T) {
	type args struct {
		openId    string
		platId    int
		worldId   int
		roleId    int
		clientIP  string
		clientVer string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"openId-01", 1, 1, 1, "127.0.0.1", "1.0.1"}, false},
		{"", args{"openId-02", 1, 1, 2, "127.0.0.1", "1.0.1"}, false},
		{"", args{"openId-03", 1, 1, 3, "127.0.0.1", "1.0.1"}, false},
		{"", args{"openId-04", 1, 1, 4, "127.0.0.1", "1.0.1"}, false},
		{"", args{"openId-05", 1, 1, 5, "127.0.0.1", "1.0.1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tssSdk.AddUser(
				tt.args.openId, tt.args.platId, tt.args.worldId, tt.args.roleId, tt.args.clientIP, tt.args.clientVer);
				(err != nil) != tt.wantErr {

				t.Errorf("TssSdk.AddUser(%v, %v, %v, %v, %v, %v) error = %v, wantErr %v",
					tt.args.openId,
					tt.args.platId,
					tt.args.worldId,
					tt.args.roleId,
					tt.args.clientIP,
					tt.args.clientVer,
					err, tt.wantErr)
			}
		})
	}
}

func TestTssSdk_DelUser(t *testing.T) {
	type args struct {
		openId  string
		platId  int
		worldId int
		roleId  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"openId-01", 1, 1, 1}, false},
		{"", args{"openId-02", 1, 1, 2}, false},
		{"", args{"openId-03", 1, 1, 3}, false},
		{"", args{"openId-04", 1, 1, 4}, false},
		{"", args{"openId-05", 1, 1, 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tssSdk.DelUser(
				tt.args.openId, tt.args.platId, tt.args.worldId, tt.args.roleId);
				(err != nil) != tt.wantErr {
				t.Errorf("TssSdk.DelUser(%v, %v, %v, %v) error = %v, wantErr %v",
					tt.args.openId, tt.args.platId, tt.args.worldId, tt.args.roleId, err, tt.wantErr)
			}
		})
	}
}

func TestTssSdk_OnRecvAntiData(t *testing.T) {
	type args struct {
		openId   string
		platId   int
		worldId  int
		roleId   int
		antiData string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"openId-01", 1, 1, 1, "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tssSdk.OnRecvAntiData(
				tt.args.openId, tt.args.platId, tt.args.worldId, tt.args.roleId, tt.args.antiData);
				(err != nil) != tt.wantErr {

				t.Errorf("TssSdk.OnRecvAntiData(%v, %v, %v, %v, %v) error = %v, wantErr %v",
					tt.args.openId, tt.args.platId, tt.args.worldId, tt.args.roleId, tt.args.antiData, err, tt.wantErr)
			}
		})
	}
}

func TestTssSdk_Tick(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tssSdk.Tick()
		})
	}
}

func TestTssSdk_Unload(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tssSdk.Unload()
		})
	}
}
