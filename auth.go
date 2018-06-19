package main

import (
	"crypto/md5"
	"net/http"
	"path/filepath"
	"text/template"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fullpath := filepath.Join(templateFolder, "login.html")
		t := template.Must(template.ParseFiles(fullpath))

		data := map[string]interface{}{
			"Host": host,
		}
		t.Execute(w, data)

	case "POST":
		r.ParseForm()

		db, err := NewDB()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		h := md5.New()
		username := r.FormValue("username")
		password := string(h.Sum([]byte(r.FormValue("password"))))

		if user, err := db.QueryUser(username, password); err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: username,
				Path:  "/",
			})

			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.Write([]byte(err.Error()))
		}
	}
}
