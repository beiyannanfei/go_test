package send_mail

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"
)

func Start() {
	for {
		now := time.Now()
		logrus.Info(now)
		// 计算下一个整小时
		next := now.Add(time.Minute * 5)
		//next := now.Add(time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		//next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0,0,0,next.Location())
		time.Sleep(next.Sub(now))
		waringMail()
	}
}

func waringMail() {
	logrus.Info("enter waring mail")
	go ordersWatch()
	go paymentWatch()
	go itemWatch()
	go deviceWatch()
	go accountWatch()
}

const ItemWatch = "item"
const DeviceWatch = "device"
const AccountWatch = "account"
const PaymentWatch = "payment"
const OrdersWatch = "orders"
const OrdersWaring = "ordersWaring"

// 订单报警
func DoOrderWaring(orderType int, sequence string, vendor model.Vendor) {
	flag := storage.UpdateOrdersState(orderType, sequence, vendor)
	if true != flag {
		// 今日已经发过报警邮件
		return
	}

	createCount := storage.GetOrderCount(storage.PreGenSequence, vendor)
	finishCount := storage.GetOrderCount(storage.QuerySequence, vendor)

	threshold := watch_waring.GetFailedOrderWaring()
	failedCount := createCount - finishCount
	if failedCount < threshold {
		return
	}

	storage.SetOrderWaringDone(vendor)
	info := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%d</td></tr>", vendor, time.Now().Format("2006-01-02 03:04:05"),
		finishCount, failedCount)
	body := `
    <html>
    	<body>
    		<h3>失败订单报警</h3>
			<table border="1">
				<tr><th>渠道名称</th><th>时间</th><th>成功订单数量</th><th>失败订单数量</th></tr>
				`+info+`
			</table>
    	</body>
    </html>
    `

	SendMail("失败订单报警", body, OrdersWaring, settings.GetString("waring", "waringSenderName"))
}

// 定时订单数据监控
func ordersWatch() {
	channelList := storage.GetOrderChannelList()
	info := ""
	//threshold := watch_waring.GetOrderWatch()
	for _, channel := range channelList {
		createCount := storage.GetOrderCount(storage.PreGenSequence, model.Vendor(channel))
		finishCount := storage.GetOrderCount(storage.QuerySequence, model.Vendor(channel))
		//if createCount - finishCount >= threshold {
			info = info + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%d</td></tr>", channel, time.Now().Format("2006-01-02 03:04:05"),
				finishCount, createCount-finishCount)
		//}
	}

	if len(info) <= 0 {
		logrus.Info("orders watch nothing need send")
		return
	}
	logrus.WithFields(logrus.Fields{
		"info": info,
	}).Info("order watch")

	body := `
    <html>
	    <body>
		    <h3>失败订单监控</h3>
			<table border="1">
				<tr><th>渠道名称</th><th>时间</th><th>成功订单数量</th><th>失败订单数量</th></tr>
				` + info + `
			</table>
	    </body>
    </html>
    `

	SendMail("失败订单监控", body, OrdersWatch, settings.GetString("waring", "waringSenderName"))
	logrus.Info("send orders watch mail success")
}

// 异常充值数据监控
func paymentWatch() {
	info := storage.GetPaymentInfo(float64(watch_waring.GetPaymentWatch()))

	if len(info) <= 0 {
		logrus.Info("payment watch nothing need send")
		return
	}
	logrus.WithFields(logrus.Fields{
		"info": info,
	}).Info("payment watch")

	body := `
    <html>
	    <body>
		    <h3>订单充值监控</h3>
			<table border="1">
				<tr><th>RoleId</th><th>时间</th><th>渠道号</th><th>当日累计充值金额</th></tr>
				` + info + `
			</table>
    	</body>
    </html>
    `
	SendMail("订单充值监控", body, PaymentWatch, settings.GetString("waring", "waringSenderName"))
	logrus.Info("send payment watch mail success")
}

// 异常物品、货币数据监控
func itemWatch() {
	info := ""
	itemMap := watch_waring.GetItemWatchMap()
	for itemId, template := range itemMap {
		info += storage.GetItemWatchInfo(itemId, 0, template.Add)
		info += storage.GetItemWatchInfo(itemId, 1, template.Sub)
	}

	if len(info) <= 0 {
		logrus.Info("item watch nothing need send")
		return
	}
	logrus.WithFields(logrus.Fields{
		"info": info,
	}).Info("item watch")

	body := `
    <html>
    	<body>
    		<h3>道具、货币数量监控</h3>
			<table border="1">
				<tr><th>RoleId</th><th>时间</th><th>渠道号</th><th>ItemId</th><th>操作类型</th><th>累计数量</th></tr>
				` + info + `
			</table>
    	</body>
    </html>
    `

	SendMail("道具、货币数量监控", body, ItemWatch, settings.GetString("waring", "waringSenderName"))
	logrus.Info("send item watch mail success")
}

// 异常设备数据监控
func deviceWatch() {
	info := storage.GetDeviceWatchInfo(watch_waring.GetDeviceWatch())

	if len(info) <= 0 {
		logrus.Info("device watch nothing need send")
		return
	}
	logrus.WithFields(logrus.Fields{
		"info": info,
	}).Info("device watch")

	body := `
    <html>
	    <body>
		    <h3>异常设备监控</h3>
			<table border="1">
				<tr><th>DeviceId</th><th>时间</th><th>角色数量</th></tr>
				` + info + `
			</table>
    	</body>
    </html>
`
	SendMail("异常设备监控", body, DeviceWatch, settings.GetString("waring", "waringSenderName"))
	logrus.Info("send device watch mail success")
}

// 异常账号数据监控
func accountWatch() {
	info := storage.GetAccountWatchInfo(watch_waring.GetAccountWatch())

	if len(info) <= 0 {
		logrus.Info("account watch nothing need send")
		return
	}
	logrus.WithFields(logrus.Fields{
		"info": info,
	}).Info("account watch")

	body := `
    <html>
	    <body>
		    <h3>异常账号监控</h3>
			<table border="1">
				<tr><th>RoleId</th><th>时间</th><th>渠道号</th><th>设备数量</th></tr>
				` + info + `
			</table>
    	</body>
    </html>
`
	SendMail("异常账号监控", body, AccountWatch, settings.GetString("waring", "waringSenderName"))
	logrus.Info("send account watch mail success")
}

func SendMail(subject, body string, toGroup string, senderName string) {
	host := settings.GetString("waring", "host")
	port := settings.GetInt("waring", "port")
	email := settings.GetString("waring", "sender")
	password := settings.GetString("waring", "password")
	var toEmail = settings.GetStringSlice("waring", toGroup+"ToEmail")
	//logrus.Info(toEmail)
	if len(toEmail) <= 0 {
		logrus.Info("no mail receivers here, groupName = ", toGroup)
		return
	}

	header := make(map[string]string)
	header["From"] = senderName + "<" + email + ">"
	//logrus.Info(header)
	header["To"] = strings.Join(toEmail, ";")
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body


	logrus.WithFields(logrus.Fields{
		"message": message,
	}).Info("waring_watch mail")
	auth := LoginAuth(email, password)	// 兼容go 1.92版本加密认证
	//auth := smtp.PlainAuth(
	//	"",
	//	email,
	//	password,
	//	host,
	//)
	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		toEmail,
		[]byte(message),
	)
	if err != nil {
		logrus.Error(err.Error())
		panic(err)
	} else {
		fmt.Println("Send watch_waring success!")
	}
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		logrus.Info(addr)
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
