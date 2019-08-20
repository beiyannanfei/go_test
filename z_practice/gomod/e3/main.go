package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"fmt"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("start...")
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
