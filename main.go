package main

import (
	"net/http"
	"os"
)

const (
	templateFolder = "templates"
)

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/game", gameHandler)

	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
