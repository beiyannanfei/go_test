package mangguotv_sdk

import (
	"fmt"
	"testing"
)

func TestValidateUser(t *testing.T) {

	err := ValidateUser("mgyxwx_KPyogdMNJxH1BDbFms6I", "9faf9f8981334a5290e73974c95d85f1")
	fmt.Println(err)
}
