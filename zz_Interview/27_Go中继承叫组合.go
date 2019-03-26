package main

import "fmt"

type student27 struct {
	Name string
	Age  int
}

func (p *student27) love() {
	fmt.Println("student27 love")
}

func (p *student27) like() {
	fmt.Println("student27 like first")
	p.love()
}

type body27 struct {
	student27
}

func (b *body27) love() {
	fmt.Println("body27 love")
}

func main() {
	b := body27{student27{"aaa", 20}}
	b.love() //out: body27 love
	b.like() //out: student27 like first student27 love
	fmt.Println(b.Age)
	fmt.Println(b.Name)
}
