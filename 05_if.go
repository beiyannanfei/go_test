/**
 * Created by wyq on 18/2/4.
 */
package main

import "fmt"

func if_test() {
	var a = 20;
	if (a > 10) {
		fmt.Println("a > 10");
		return
	}
	fmt.Println("a <= 10");
	return
}

func if_test1() {
	var a = false
	if (!!a) {
		fmt.Println("a is true")
		return
	}
	fmt.Println("a is false")
	return
}

func main() {
	if_test()
	if_test1()
}
