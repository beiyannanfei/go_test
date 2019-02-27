package main

import (
	"net/url"
	"fmt"
)

func main() {
	request := "data=%7B%22client_ip%22%3A%22223.104.3.183%22%2C%22extension_info%22%3A%221833157522294.3.2019-01-14T23%3A12%3A00%2B08%3A00.37017535139084%22%2C%22game_id%22%3A%221165%22%2C%22game_money%22%3A%2210%22%2C%22id%22%3A%2232906315%22%2C%22merchant_id%22%3A%22599%22%2C%22money%22%3A%22100%22%2C%22order_no%22%3A%225474787221773392%22%2C%22order_status%22%3A1%2C%22out_trade_no%22%3A%221833157522294-1547478720%22%2C%22pay_money%22%3A%22100%22%2C%22pay_time%22%3A%221547478761%22%2C%22product_desc%22%3A%22%5Cu8863%5Cu4e4b%5Cu613f%5Cu793c%5Cu5305%22%2C%22product_name%22%3A%22%5Cu8863%5Cu4e4b%5Cu613f%5Cu793c%5Cu5305%22%2C%22role%22%3A%22sxmad%22%2C%22sign%22%3A%2246b1c1208d4f8bb620df1e992624ed68%22%2C%22uid%22%3A%22396481564%22%2C%22username%22%3A%22sxmad%22%2C%22zone_id%22%3A%221386%22%7D"
	a, err := url.ParseQuery(request)
	fmt.Printf("%#v\n", a)
	fmt.Println(err)
	fmt.Println(a.Get("data"))
	//v := url.Values{}
	//v.Add("msg", "此订单不存在或已经提交")
	//body := v.Encode()
	//fmt.Println(v)
	//fmt.Println(body)
	// url decode
	//m, _ := url.ParseQuery(body)
	//fmt.Println(m)
}
