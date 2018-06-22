package main

import (
	"errors"
	"github.com/avalchev94/office_worldcup/database"
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

func authenticated(w http.ResponseWriter, r *http.Request) (database.User, error) {
	c, err := r.Cookie("auth")
	if err != nil || c.Value == "" {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return database.User{}, errors.New("auth cookie not found")
	}

	m := objx.MustFromBase64(c.Value)
	return database.User{
		ID:        database.ObjectIdHex(m["id"].(string)),
		Username:  m["username"].(string),
		Favourite: database.ObjectIdHex(m["favourite"].(string)),
		Fullname:  m["fullname"].(string),
		Role:      m["role"].(string),
		Points:    int32(m["points"].(int)),
	}, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	switch r.Method {
	case "GET":
		fullpath := filepath.Join(templateFolder, "login.html")
		t := template.Must(template.ParseFiles(fullpath))

		t.Execute(w, nil)

	case "POST":
		r.ParseForm()

		db, err := database.New()
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
			"id":        user.ID.Hex(),
			"username":  user.Username,
			"fullname":  user.Fullname,
			"favourite": user.Favourite.Hex(),
			"points":    user.Points,
			"role":      user.Role,
			"avatar":    user.Avatar,
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
	if _, err := r.Cookie("auth"); err == nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	switch r.Method {
	case "GET":
		fullpath := filepath.Join(templateFolder, "register.html")
		t := template.Must(template.ParseFiles(fullpath))

		db, err := database.New()
		defer db.Close()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		dbTeams, _ := db.GetTeams()
		data := map[string]interface{}{
			"Teams": make([]map[string]string, 0),
		}

		for _, t := range dbTeams {
			teams := data["Teams"].([]map[string]string)
			data["Teams"] = append(teams, map[string]string{
				"ID":   t.ID.Hex(),
				"Name": t.Name,
			})
		}

		t.Execute(w, data)
	case "POST":
		r.ParseForm()

		user := database.User{
			Username:  r.FormValue("username"),
			Fullname:  r.FormValue("fullname"),
			Favourite: database.ObjectIdHex(r.FormValue("favourite")),
			Password:  encryptPassword(r.FormValue("password")),
			Role:      "user",
			Points:    0,
			Avatar:    "",
		}

		db, err := database.New()
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
