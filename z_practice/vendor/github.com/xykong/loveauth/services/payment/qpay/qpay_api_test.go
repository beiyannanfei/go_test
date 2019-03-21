package qpay

import (
	"testing"
	"fmt"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

func TestUnifiedorder(t *testing.T) {
	order := model.Order{}
	order.GlobalId = 558877
	order.ShopId = 12314
	order.Sequence = "2362765424644.2.2018-05-30T10:10:59+08:00.2363086"
	order.Amount = 10
	order.Timestamp = time.Now()

	test := struct {
		name  string
		ip    string
		order *model.Order
	}{
		"TestUnifiedorder",
		"127.0.0.1",
		&order,
	}
	t.Run(test.name, func(t *testing.T) {
		resp, err := Unifiedorder(test.ip, test.order)
		fmt.Println(resp, err)
		fmt.Printf("%#v\n", resp)
	})
}
