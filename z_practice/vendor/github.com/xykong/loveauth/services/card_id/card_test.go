package card_id

import (
	"testing"
	"fmt"
)

func Test_IsAdult(t *testing.T) {
	type args struct {
		id string
	}

	test := struct {
		name string
		args args
	}{"isAdult", args{"131022198807181135"}}

	t.Run(test.name, func(t *testing.T) {
		res := IsAdult(test.args.id)
		fmt.Println(res)
	})
}
