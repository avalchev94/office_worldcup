package main

import (
	"github.com/avalchev94/office_worldcup/database"
	"net/http"
	"regexp"
)

func betHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if user, err := authenticated(w, r); err == nil {
			r.ParseForm()

			db, err := database.New()
			defer db.Close()
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			for key, value := range r.Form {
				if _, err := db.GetMatch(database.ObjectIdHex(key)); err == nil {
					r, _ := regexp.Compile("[0-9]+:[0-9]+")
					if !r.MatchString(value[0]) {
						w.Write([]byte("Incorrect input"))
						return
					}

					p := database.Prediction{
						Match:     database.ObjectIdHex(key),
						User:      user.ID,
						Predicted: value[0],
						Score:     0,
					}

					if err = db.AddPrediction(p); err != nil {
						w.Write([]byte(err.Error()))
						return
					}
				}
			}

			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
