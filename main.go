package main

import (
	"github.com/avalchev94/office_worldcup/admin"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

const (
	templateFolder = "templates"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/bet", betHandler)
	r.HandleFunc("/highscore", highscoreHandler)
	r.HandleFunc("/predictions", predictionsHandler)

	admin.Handle(r.PathPrefix("/admin").Subrouter())

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
