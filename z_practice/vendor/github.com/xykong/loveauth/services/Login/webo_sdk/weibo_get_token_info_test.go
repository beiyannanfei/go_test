package webo_sdk

import (
	"testing"
	"fmt"
)

func TestWeiboGetTokenInfo(t *testing.T) {
	test := struct {
		name         string
		access_token string
	}{
		"TestWeiboGetTokenInfo", "2.00BBZKFCKQ7_HCdb757c317eKedo2D",
	}

	t.Run(test.name, func(t *testing.T) {
		resp, err := WeiBoGetTokenInfo(test.access_token)
		fmt.Println(resp, err)
	})
}
