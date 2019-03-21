package dict

import (
	"fmt"
	"testing"
)

func TestExistInvalidWord(t *testing.T) {

	Load("../../templates/maskingwords.txt")
	a, b := ExistInvalidWord("习近平是个江泽民")
	c, d := ExistInvalidWord("你好啊大兄弟")
	e, f, g := ReplaceInvalidWords("习近平是个江泽民")
	fmt.Println(a, b)
	fmt.Println(c, d)
	fmt.Println(e, f, g)
}
