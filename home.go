package main

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if user, err := authenticated(w, r); err == nil {
		db, err := NewDB()
		defer db.Close()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		matchesDB, err := db.GetTodayMatches()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		data := map[string]interface{}{
			"Host": host,
		}

		data["Matches"] = make([]map[string]string, 0)
		for _, m := range matchesDB {
			if _, err := db.GetPrediction(m.ID, user.ID); err != nil {
				host, _ := db.GetTeamByID(m.Host)
				guest, _ := db.GetTeamByID(m.Guest)

				matches := data["Matches"].([]map[string]string)
				data["Matches"] = append(matches, map[string]string{
					"ID":    m.ID.Hex(),
					"Host":  host.Name,
					"Guest": guest.Name,
					"Date":  m.Date.Format("15:04"),
				})
			}
		}

		fullpath := filepath.Join(templateFolder, "home.html")
		t := template.Must(template.ParseFiles(fullpath))
		t.Execute(w, data)
	}
}
