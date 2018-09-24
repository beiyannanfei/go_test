package main

import (
	"net/http"
	"fmt"
	"strings"
)

type myMux struct {
}

func (p *myMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ServeHTTP Path: %v\n", r.URL.Path)
	if r.URL.Path == "/" {
		sayHelloName1(w, r)
		return
	}

	http.NotFound(w, r)
	return
}

func sayHelloName1(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path:", r.URL.Path)
	fmt.Println("scheme:", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, " "))
	}
	fmt.Fprint(w, "Hello myroute!") //这个写入到w的是输出到客户端的
}

func main() {
	mux := &myMux{}
	http.ListenAndServe(":9090", mux)
}

//curl "http://localhost:9090/?url_long=111&url_long=222"
//curl "http://localhost:9090"