package main

import (
	"errors"
	"github.com/stretchr/objx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"path/filepath"
	"text/template"
)

func encryptPassword(password string) string {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encrypted)
}

func comparePasswords(encrypted, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password)) == nil
}

func authenticated(w http.ResponseWriter, r *http.Request) (User, error) {
	c, err := r.Cookie("auth")
	if err != nil || c.Value == "" {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return User{}, errors.New("auth cookie not found")
	}

	m := objx.MustFromBase64(c.Value)
	return User{
		Username: m["username"].(string),
		Fullname: m["fullname"].(string),
		Role:     m["role"].(string),
		Points:   int32(m["points"].(int)),
	}, nil
}

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
		defer db.Close()

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		user, err := db.GetUser(r.FormValue("username"))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if !comparePasswords(user.Password, r.FormValue("password")) {
			w.Write([]byte("Password is incorrect."))
			return
		}

		authCookieVale := objx.New(map[string]interface{}{
			"username": user.Username,
			"fullname": user.Fullname,
			"points":   user.Points,
			"role":     user.Role,
			"avatar":   user.Avatar,
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieVale,
			Path:  "/",
		})

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fullpath := filepath.Join(templateFolder, "register.html")
		t := template.Must(template.ParseFiles(fullpath))

		data := map[string]interface{}{
			"Host": host,
		}
		t.Execute(w, data)
	case "POST":
		r.ParseForm()

		user := User{
			Username: r.FormValue("username"),
			Fullname: r.FormValue("fullname"),
			Password: encryptPassword(r.FormValue("password")),
			Role:     "user",
			Points:   0,
			Avatar:   "",
		}

		db, err := NewDB()
		defer db.Close()

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if err := db.AddUser(user); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "login", http.StatusMovedPermanently)
	}
}
