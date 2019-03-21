package tss

import (
	"net"
	"fmt"
	"bufio"
	"github.com/xykong/loveauth/settings"
	"encoding/hex"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var quitSemaphore chan bool
var conn *net.TCPConn

func Start() {

	if !settings.GetBool("tencent", "tss.enable") {
		logrus.WithFields(logrus.Fields{
			"enable": settings.GetBool("tencent", "tss.enable"),
		}).Warn("TSS is disabled.")
		return
	}

	var address = settings.GetString("tencent", "tss.address")

	if len(address) == 0 {
		logrus.WithFields(logrus.Fields{
			"address": address,
		}).Fatal("TSS Address is not valid.")
		return
	}

	var err error
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"address": address,
		}).Fatal("TSS ResolveTCPAddr error.")
		return
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"address": address,
		}).Fatal("TSS DialTCP error.")
		return
	}

	logrus.WithFields(logrus.Fields{
		"address": address,
	}).Info("TSS Services started.")

	go onMessageReaderReceived(conn)

	//var index int32 = 0
	//
	//for n := 0; n < 8; n++ {
	//	go func() {
	//		for {
	//
	//			index += 1
	//
	//			gamePkg := &GamePkg{
	//				Head: &GamePkgHead{
	//					CmdId:  int32(GameCmdID_GAME_CMDID_LOGIN_CHANNEL),
	//					Openid: fmt.Sprintf("TEST-TSS-OPENID-%v", index),
	//				},
	//				Body: &GamePkgBody{
	//					Login: &LoginChannel{
	//						AuthSignature: index,
	//						ClientVersion: 1357,
	//					},
	//					Logout:               nil,
	//					TransAntiData:        nil,
	//					TransAntiDecryptData: nil,
	//					RoleList:             nil,
	//					SelectRole:           nil,
	//				},
	//			}
	//
	//			sendPackage(conn, gamePkg)
	//		}
	//	}()
	//}

	//<-quitSemaphore
}

func sendPackage(conn *net.TCPConn, gamePkg *GamePkg) {

	data, err := proto.Marshal(gamePkg)
	if err != nil {
		logrus.Fatal("marshaling error: ", err)
	}

	bs := make([]byte, 4)

	binary.LittleEndian.PutUint32(bs, uint32(len(data)))

	//encodedStr := hex.EncodeToString(bs)
	//encodedStr += hex.EncodeToString(data)
	//fmt.Printf("%s\n", encodedStr)

	//for n := 3; n > 0; n -= 1 {

	//logrus.Info("send message.")
	conn.Write(bs)
	conn.Write(data)
}

func onMessageSimpleReceived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		fmt.Println(msg)
		if err != nil {
			quitSemaphore <- true
			break
		}
	}
}

func onMessageReaderReceived(conn *net.TCPConn) {

	//fmt.Printf("onMessageReaderReceived")

	reader := bufio.NewReader(conn)

	for {
		_, err := reader.Peek(4)
		if err != nil {
			logrus.Info("reader.Peek.")
			continue
		}

		var buf = make([]byte, 4)
		reader.Read(buf)

		length := int(binary.LittleEndian.Uint32(buf))
		if length <= 0 {
			logrus.Info("length is 0.")
			continue
		}

		_, err = reader.Peek(length)
		if err != nil {
			logrus.Info("reader.Peek.")
			continue
		}

		buf = make([]byte, length)
		reader.Read(buf)

		gamePkg := GamePkg{}
		err = proto.Unmarshal(buf, &gamePkg)
		if err != nil {
			logrus.Errorf("proto.Unmarshal failed: %v", err)
			continue
		}

		//fmt.Printf("recv: %v", gamePkg)
	}
}

func onMessageReceived(conn *net.TCPConn) {

	fmt.Printf("onMessageReceived")

	reader := bufio.NewReader(conn)

	for {

		//if reader.Buffered() < 4 {
		//	continue
		//}

		buf, err := reader.Peek(4)
		if err != nil {
			logrus.Info("reader.Peek.")
			continue
		}

		length := int(binary.LittleEndian.Uint32(buf))
		if length <= 0 {
			logrus.Info("length is 0.")
			continue
		}

		fmt.Printf("%s\n", hex.EncodeToString(buf))

		buffered := reader.Buffered()
		if length+4 > buffered {
			logrus.Infof("length is not enough: %v, %v", length, buffered)
			continue
		}

		var data = make([]byte, length+4)
		n, err := reader.Read(data[:])
		if err != nil {
			logrus.Errorf("reader.Read failed: %v", err)
			continue
		}

		if n != length+4 {
			logrus.Errorf("reader.Read failed: length is not equal: %v", n)
			continue
		}

		gamePkg := GamePkg{}
		err = proto.Unmarshal(data[4:], &gamePkg)
		if err != nil {
			logrus.Errorf("proto.Unmarshal failed: %v", err)
			continue
		}

		fmt.Printf("recv: %v", gamePkg)
	}
}

func sendMessage(gamePkg *GamePkg) {

	var address = settings.GetString("tencent", "tss.address")
	if len(address) == 0 {
		return
	}

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", address)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logrus.Fatal("DialTCP error: ", err)
	}

	defer conn.Close()

	fmt.Println("connected!")

	go onMessageReaderReceived(conn)

	data, err := proto.Marshal(gamePkg)
	if err != nil {
		logrus.Fatal("marshaling error: ", err)
	}

	bs := make([]byte, 4)

	binary.LittleEndian.PutUint32(bs, uint32(len(data)))

	encodedStr := hex.EncodeToString(bs)
	encodedStr += hex.EncodeToString(data)
	fmt.Printf("%s\n", encodedStr)

	//for n := 3; n > 0; n -= 1 {

	logrus.Info("send message.")
	conn.Write(bs)
	conn.Write(data)

	//time.Sleep(time.Second)
	//}

	//<-quitSemaphore
}

func AddUser(gamePkg *GamePkg) {

}
