package main

import (
	"net/http"
)

const (
	templateFolder = "templates"
	host           = ":1914"
)

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/game", gameHandler)

	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(host, nil)
}
