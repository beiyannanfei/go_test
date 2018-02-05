package main

import "fmt"

func main() {
	s1();
	fmt.Println("============================")

}

func structParam(book Books) { //结构体作为函数参数
	fmt.Printf("Book title : %s\n", book.title);
	fmt.Printf("Book author : %s\n", book.author);
	fmt.Printf("Book subject : %s\n", book.subject);
	fmt.Printf("Book book_id : %d\n", book.book_id);
}

type Books struct {
	title   string
	author  string
	subject string
	book_id int
}

func s1() {
	var Book1 Books
	var Book2 Books

	Book1.title = "Go 语言"
	Book1.author = "www.runoob.com"
	Book1.subject = "Go 语言教程"
	Book1.book_id = 6495407

	Book2 = Books{"Python 教程", "www.runoob.com", "Python 语言教程", 6495700}

	/* 打印 Book1 信息 */
	fmt.Printf("Book 1 title : %s\n", Book1.title)
	fmt.Printf("Book 1 author : %s\n", Book1.author)
	fmt.Printf("Book 1 subject : %s\n", Book1.subject)
	fmt.Printf("Book 1 book_id : %d\n", Book1.book_id)

	/* 打印 Book2 信息 */
	fmt.Printf("Book 2 title : %s\n", Book2.title)
	fmt.Printf("Book 2 author : %s\n", Book2.author)
	fmt.Printf("Book 2 subject : %s\n", Book2.subject)
	fmt.Printf("Book 2 book_id : %d\n", Book2.book_id)
	fmt.Println("============================")

	structParam(Book1)
	fmt.Println("============================")

	printBook(&Book1)
	printBook(&Book2)
}

func printBook(book *Books) {
	fmt.Printf("Book title : %s\n", book.title);
	fmt.Printf("Book author : %s\n", book.author);
	fmt.Printf("Book subject : %s\n", book.subject);
	fmt.Printf("Book book_id : %d\n", book.book_id);
}
