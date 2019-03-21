package douyin

import (
	"fmt"
	"testing"
)

func TestOrderStatus(t *testing.T) {

	err := OrderStatus("tta080138fc6d476", "ef7e4d852cbeab9e3669d4f103fbc5ff", "734904904382_3_20190129120835.734907617323")
	fmt.Println(err)
}
