package main

//Go 中不允许不同类型之间的混合使用，但是对于常量的类型限制非常少，因此允许常量之间的混合使用，下面这个程序很好地解释了这个现象（该程序无法通过编译）：

func main() {
	var a int
	var b int32
	a = 15
	//b = a + a // 编译错误 cannot use a + a (type int) as type int32 in assignment
	b = int32(a) + int32(a)
	b = int32(a + a)
	b = b + 5 // 因为 5 是常量，所以可以通过编译
}
