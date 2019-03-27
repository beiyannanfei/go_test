package main

import (
	"io"
	"fmt"
	"bytes"
)

const debug = true

func main() {
	var buf1, buf2 *bytes.Buffer
	if debug {
		buf1 = new(bytes.Buffer)
	}
	if !debug {
		buf2 = new(bytes.Buffer)
	}

	f(buf1)
	f(buf2)

	if buf2 == nil {
		fmt.Println("right")	//输出: right
	}
}
//todo 接口指针 传入函数的接口参数时，才会出现以上的坑
func f(out io.Writer) {
	if out != nil {
		fmt.Println("surprise!")
	}
}

//surprise!
//surprise!

//todo 这就牵扯到一个概念了，是关于接口值的。概念上讲一个接口的值分为两部分：一部分是类型，一部分是类型对应的值，他们分别叫：动态类型和动态值。类型系统是针对编译型语言的，类型是编译期的概念，因此类型不是一个值。在上述代码中，给f函数的out参数赋了一个*bytes.Buffer的空指针，所以out的动态值是nil。然而它的动态类型是bytes.Buffer，意思是：“A non-nil interface containing a nil pointer”，所以“out!=nil”的结果依然是true。
