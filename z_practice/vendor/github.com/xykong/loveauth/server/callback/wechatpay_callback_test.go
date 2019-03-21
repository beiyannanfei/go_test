package callback

import (
	"testing"
	"encoding/xml"
	"github.com/xykong/loveauth/storage/model"
	"encoding/json"
	"fmt"
	"github.com/xykong/loveauth/services/payment/wechat_pay"
)

func TestUnifiedorder(t *testing.T) {
	requestXml := "<xml><appid><![CDATA[wx67722a94f90b26a4]]></appid>\n<attach><![CDATA[2362765424644.2.2018-05-30T10:10:59+08:00.2363086403199]]></attach>\n<bank_type><![CDATA[CFT]]></bank_type>\n<cash_fee><![CDATA[40]]></cash_fee>\n<fee_type><![CDATA[CNY]]></fee_type>\n<is_subscribe><![CDATA[N]]></is_subscribe>\n<mch_id><![CDATA[1521242291]]></mch_id>\n<nonce_str><![CDATA[ickli1ojtfeSUok04ZBMv0FhjphT0PsJ]]></nonce_str>\n<openid><![CDATA[oLRcH1uhmp0dDzTbOAXJgq733u0M]]></openid>\n<out_trade_no><![CDATA[4455778-1547814634]]></out_trade_no>\n<result_code><![CDATA[SUCCESS]]></result_code>\n<return_code><![CDATA[SUCCESS]]></return_code>\n<sign><![CDATA[DC89E28490EBF6A22BED2CFD45A5D2B9]]></sign>\n<time_end><![CDATA[20190118203317]]></time_end>\n<total_fee>40</total_fee>\n<trade_type><![CDATA[APP]]></trade_type>\n<transaction_id><![CDATA[4200000260201901181628250870]]></transaction_id>\n</xml>"
	var request model.WechatPayCallback
	err := xml.Unmarshal([]byte(requestXml), &request)
	fmt.Println(err)
	jsonStr, _ := json.Marshal(request)
	fmt.Println(string(jsonStr))

	sign := wechat_pay.GeneratePaySign(string(jsonStr))
	fmt.Println("sign:", sign)
}
