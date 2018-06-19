package main

import (
	"net/http"
)

const (
	templateFolder = "templates"
	host           = ":1914"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/login", loginHandler)
	//	http.HandleFunc("/register", registerHandler)

	http.ListenAndServe(host, nil)
}
