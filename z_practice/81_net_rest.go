package main

import (
	"net/http"
	"fmt"
	"github.com/drone/routes"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	uid := params.Get(":uid")
	fmt.Fprintf(w, "you are get user %s", uid)
}

func modifyUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	uid := params.Get(":uid")
	fmt.Fprint(w, "you are modify user: %s", uid)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	uid := params.Get(":uid")
	fmt.Fprint(w, "you are delete user %s", uid)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	uid := params.Get(":uid")
	fmt.Fprint(w, "you are add user %s", uid)
}

func main() {
	mux := routes.New()
	mux.Get("/user/:uid", getUser)			//curl "127.0.0.1:8088/user/1234"
	mux.Post("/user/:uid", modifyUser)		//curl -X POST "127.0.0.1:8088/user/1234"
	mux.Del("/user/:uid", deleteUser)		//curl -X DELETE "127.0.0.1:8088/user/1234"
	mux.Put("/user/:uid", addUser)			//curl -X PUT "127.0.0.1:8088/user/1234"
	http.Handle("/", mux)
	http.ListenAndServe(":8088", nil)
}
