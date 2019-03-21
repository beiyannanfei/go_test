package mangguotv

import (
	"fmt"
	"net/url"
	"testing"
)

func TestCheckSign(t *testing.T) {

	value := url.Values{}
	value.Add("uuid", "112")
	value.Add("appId", "1212")
	value.Add("sign","4bc7c2ed21b95e2c1643b0e612679b80e34c67d3")

	result :=CheckSign(value, "xMVSVy9kZc1sDlLs")
	fmt.Println(result)
}
