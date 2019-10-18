package main

import "fmt"

type People struct{}

type Teacher struct {
	People
}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}

func (p *People) ShowB() {
	fmt.Println("showB")
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func main() {
	t := Teacher{}
	t.ShowA() // => t.People.ShowA()
	t.ShowB()
}

/*
showA
showB
teacher showB
*/
