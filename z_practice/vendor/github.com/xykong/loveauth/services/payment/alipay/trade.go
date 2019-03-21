package alipay

import (
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"net/url"
)

// TradePagePay https://docs.open.alipay.com/270/alipay.trade.page.pay
func TradePagePay(param AliPayTradePagePay) (results *url.URL, err error) {

	request := settings.GetString("lovepay", "alipay.url")

	p, err := URLValues(param)
	if err != nil {

		return nil, err
	}

	results, err = url.Parse(request + "?" + p.Encode())
	if err != nil {

		return nil, err
	}

	return results, err
}

// TradeAppPay https://docs.open.alipay.com/api_1/alipay.trade.app.pay
func TradeAppPay(param AliPayTradeAppPay) (results string, err error) {

	p, err := URLValues(param)
	if err != nil {

		return "", err
	}

	return p.Encode(), err
}

// TradePreCreate https://docs.open.alipay.com/api_1/alipay.trade.precreate/
func TradePreCreate(param AliPayTradePreCreate) (results *AliPayTradePreCreateResponse, err error) {

	request := settings.GetString("lovepay", "alipay.url")

	err = doRequest(request, param, &results)

	return results, err
}

// TradeQuery https://docs.open.alipay.com/api_1/alipay.trade.query/
func TradeQuery(param AliPayTradeQuery) (results *AliPayTradeQueryResponse, err error) {

	request := settings.GetString("lovepay", "alipay.url")

	err = doRequest(request, param, &results)

	if results.AliPayTradeQuery.Code != "10000" {

		return results, errors.New("query trade code err :" + results.AliPayTradeQuery.Code + " msg:" + results.AliPayTradeQuery.Msg)
	}

	if !(results.AliPayTradeQuery.TradeStatus == K_TRADE_STATUS_TRADE_SUCCESS ||
		results.AliPayTradeQuery.TradeStatus == K_TRADE_STATUS_TRADE_FINISHED) {

		return results, errors.New("query trade status err :" + results.AliPayTradeQuery.TradeStatus)
	}

	return results, err
}

// TradeRefund https://docs.open.alipay.com/api_1/alipay.trade.refund/
func TradeRefund(param AliPayTradeRefund) (results *AliPayTradeRefundResponse, err error) {

	request := settings.GetString("lovepay", "alipay.url")

	err = doRequest(request, param, &results)

	return results, err
}
