package foo

import (
	"testing"
	"time"
	"fmt"
)

//Go的单元测试函数分为两类:功能测试函数和性能测试函数，分别为以Test和Benchmark 6 为函数名前缀并以*testing.T为单一参数的函数

func TestAdd1(t *testing.T) { //功能测试函数以Test为函数名前缀并以*testing.T为单一参数的函数
	r := Add(1, 2)
	if r != 3 {
		t.Errorf("Add(1, 2) failed. Got %d, expected 3.", r)
	}
}

func Add(a, b int) int {
	return a + b
}

//执行功能测试
//cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/01_GO语言编程/02_测试 && go test go_test.go

func delay(n int) {
	fmt.Printf("=============== 开始延时%d秒\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Printf("=============== 结束延时%d秒\n", n)
}

func BenchmarkAdd1(b *testing.B) { //性能测试函数以Benchmark 6 为函数名前缀并以*testing.B为单一参数的函数
	//可以看出，性能测试与功能测试代码相比，最大的区别在于代码里的这个for循环，循环b.N次。
	// 写这个for循环的原因是为了能够让测试运行足够长的时间便于进行平均运行时间的计算。 
	// 如果测试代码中一些准备工作的时间太长，我们也可以这样处理以明确排除这些准备工作所花费 时间对于性能测试的时间影响:

	b.StopTimer()  //暂停计时器
	delay(1)       //一个耗时较长的准备工作，比如读文件
	b.StartTimer() //开启计时器，之前的准备时间未计入总花费时间内

	for i := 0; i < b.N; i++ {
		result := Add(i, i)
		fmt.Printf("i: %d, result: %d\n", i, result)
		//Add(i, i)
	}
}

//执行功能测试
//cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/01_GO语言编程/02_测试 && go test -bench=. -benchtime="3s" -timeout="5s" -benchmem -cpu=8
//性能测试参数参考文档： https://www.jianshu.com/p/1adc69468b6f
