package douyin_sdk

import (
	"fmt"
	"testing"
)

func TestCheckUser(t *testing.T) {

	err := CheckUser("tta080138fc6d476", "ef7e4d852cbeab9e3669d4f103fbc5ff", "29eac12-4036-77a2-89ab-53f07acb0080d")
	fmt.Println(err)
}
