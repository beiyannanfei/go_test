package bilibili

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAuthToken(t *testing.T) {
	resp, err := AuthToken("a79b15f60aa47d2c6ed272d93cd63911", 396063952)
	data, _ := json.Marshal(resp)
	fmt.Println(string(data), err)
}
