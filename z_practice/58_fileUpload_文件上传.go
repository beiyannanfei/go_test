package main

import (
	"net/http"
	"fmt"
	"html/template"
	"os"
	"io"
)

func fileSel(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("method: %v, path: %v\n", r.Method, r.URL.Path)
	t, _ := template.ParseFiles("/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/58_fileUpload.html")
	t.Execute(w, nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("method: %v, path: %v\n", r.Method, r.URL.Path)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println("read uploadfile err:", err)
		return
	}

	defer file.Close()
	fmt.Fprint(w, "%v", handler.Header)
	fmt.Printf("file: %v, handler: %v\n", file, handler)

	f, err := os.OpenFile("/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/58_z"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println("open file err:", err)
		return
	}

	defer f.Close()
	io.Copy(f, file)
	fmt.Printf("file upload finish fileName: %v, size: %v\n", handler.Filename, handler.Size)
}

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/fileSel", fileSel)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("ListenAndServe err:", err)
	}
}
