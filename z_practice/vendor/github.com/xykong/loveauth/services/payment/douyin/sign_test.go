package douyin

import (
	"fmt"
	"net/url"
	"testing"
)

func TestCheckSign(t *testing.T) {

	value, err := url.ParseQuery("total_fee=100&client_id=tta080138fc6d476&pay_time=2019-01-29+12%3A09%3A36&trade_status=0&trade_no=2019012922001491241014433046&tt_sign=ZbUPwP5%2Bgeh0B1n2RfU958prokjYNTy%2BCv23huWcWctagzhlH93w57gtmuDi%2FDlYYrR7i4oU6FnRa7UFFS%2Bf1wehV2p%2FoHrriI0hmsYx996TCZcWyAf%2Fae6%2BJT9plfaL9eMwC9NCzUQdfpfw6HhsAEqdaA6QecIuuvYiqB5aEgg%3D&out_trade_no=734904904382_3_20190129120835.734907617323&way=2&notify_time=2019-01-29+12%3A09%3A37&tt_sign_type=RSA&notify_id=2019012900222120936091241040238319&notify_type=trade_status_sync&buyer_id=13552612360")
	if err != nil {

		fmt.Println(err)
	}

	payKey := "9525d16ac07a780b9d1ed18c15ec1394"
	fmt.Println(CheckSign(value, payKey))
}
