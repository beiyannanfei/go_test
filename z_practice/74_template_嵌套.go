package main

import (
	"html/template"
	"os"
	"fmt"
)

func main() {
	s1, _ := template.ParseFiles(
		"/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/74_header.tmpl",
		"/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/74_content.tmpl",
		"/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/74_footer.tmpl")
	s1.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println()

	s1.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println()

	s1.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()

	s1.Execute(os.Stdout, nil) //执行s1.Execute，没有任何的输出，因为在默认的情况下没有默认的子模板，所以不会输出 任何的东西
}
