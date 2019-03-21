package tss

import (
	"testing"
	"fmt"
)

//
//func TestStart(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			//Start()
//		})
//	}
//}
//
//func Test_onMessageReceived(t *testing.T) {
//	type args struct {
//		conn *net.TCPConn
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			//onMessageReceived(tt.args.conn)
//		})
//	}
//}

func Test_sendMessage(t *testing.T) {
	type args struct {
		gamePkg *GamePkg
	}

	tests := []struct {
		name string
		args args
	}{
		{"", args{&GamePkg{
			Head: &GamePkgHead{
				CmdId:  int32(GameCmdID_GAME_CMDID_LOGIN_CHANNEL),
				Openid: "TEST-TSS-OPENID-01",
			},
			Body: &GamePkgBody{
				Login: &LoginChannel{
					AuthSignature: 1234,
					ClientVersion: 1357,
				},
				Logout:               nil,
				TransAntiData:        nil,
				TransAntiDecryptData: nil,
				RoleList:             nil,
				SelectRole:           nil,
			},
		}}},
		{"", args{&GamePkg{
			Head: &GamePkgHead{
				CmdId:  int32(GameCmdID_GAME_CMDID_LOGIN_CHANNEL),
				Openid: "TEST-TSS-OPENID-02",
			},
			Body: &GamePkgBody{
				Login: &LoginChannel{
					AuthSignature: 1234,
					ClientVersion: 1357,
				},
				Logout:               nil,
				TransAntiData:        nil,
				TransAntiDecryptData: nil,
				RoleList:             nil,
				SelectRole:           nil,
			},
		}}},
		{"", args{&GamePkg{
			Head: &GamePkgHead{
				CmdId:  int32(GameCmdID_GAME_CMDID_LOGIN_CHANNEL),
				Openid: "TEST-TSS-OPENID-03",
			},
			Body: &GamePkgBody{
				Login: &LoginChannel{
					AuthSignature: 1234,
					ClientVersion: 1357,
				},
				Logout:               nil,
				TransAntiData:        nil,
				TransAntiDecryptData: nil,
				RoleList:             nil,
				SelectRole:           nil,
			},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//sendMessage(tt.args.gamePkg)
		})
	}
}

func Benchmark_sendPackage(b *testing.B) {
	Start()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		//fmt.Print("send package.")

		gamePkg := &GamePkg{
			Head: &GamePkgHead{
				CmdId:  int32(GameCmdID_GAME_CMDID_LOGIN_CHANNEL),
				Openid: fmt.Sprintf("TEST-TSS-OPENID-%v", i),
			},
			Body: &GamePkgBody{
				Login: &LoginChannel{
					AuthSignature: int32(i),
					ClientVersion: 1357,
				},
				Logout:               nil,
				TransAntiData:        nil,
				TransAntiDecryptData: nil,
				RoleList:             nil,
				SelectRole:           nil,
			},
		}

		sendPackage(conn, gamePkg)
	}
}
