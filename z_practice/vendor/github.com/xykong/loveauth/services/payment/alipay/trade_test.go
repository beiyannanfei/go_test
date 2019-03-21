package alipay

import (
	"encoding/json"
	"testing"
)

func TestAliPay_TradeAppPay(t *testing.T) {
	t.Log("========== TradeAppPay ==========")
	var p = AliPayTradeAppPay{}
	p.NotifyURL = "https://alipay.api.worldoflove.cn/api/v1/alipay/callback"
	p.Body = "body"
	p.Subject = "商品标题"
	p.OutTradeNo = "01010101"
	p.TotalAmount = "100.00"
	p.ProductCode = "p_1010101"
	param, err := TradeAppPay(p)
	if err != nil {
		t.FailNow()
	}
	t.Log(param)
}

func TestAliPay_TradePagePay(t *testing.T) {

	Init()

	t.Log("========== TradePagePay ==========")
	var p = AliPayTradePagePay{}
	p.NotifyURL = "https://auth-beta.api.worldoflove.cn/api/v1/alipay/callback"
	p.ReturnURL = "http://alipay.api.worldoflove.cn"
	p.Subject = "商品标题"
	p.OutTradeNo = "trade_test_112312313131223"
	p.TotalAmount = "0.01"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, err := TradePagePay(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAliPay_TradeQuery(t *testing.T) {
	t.Log("========== TradeQuery ==========")
	var p = AliPayTradeQuery{}
	p.OutTradeNo = "trade_test_112312313131223"

	rsp, err := TradeQuery(p)
	if err != nil {
		t.Log(err)
	}

	data,_:=json.Marshal(rsp)

	t.Log(string(data))
}
