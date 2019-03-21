package vivo_sdk

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAuthToken(t *testing.T) {
	test := struct {
		name      string
		authToken string
	}{
		"TestAuthToken",
		"db2a731bcddcf06e_db2a731bcddcf06e_e43e848018b7632210cd421e5f62cd51",
	}

	t.Run(test.name, func(t *testing.T) {

		resp, err := AuthToken(test.authToken)
		j, _ := json.Marshal(resp)

		fmt.Println(string(j), err)
	})
}
