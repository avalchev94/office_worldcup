package admin

import (
	"github.com/avalchev94/office_worldcup/database"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"text/template"
)

const (
	templateFolder = "templates"
)

func Handle(router *mux.Router) {
	router.HandleFunc("", adminHandler)
	router.HandleFunc("/add", addGameHandler)
	router.HandleFunc("/finish/{id}", finishGameHandler)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fullpath := filepath.Join(templateFolder, "admin.html")
	t := template.Must(template.ParseFiles(fullpath))

	db, _ := database.New()
	defer db.Close()

	data := map[string]interface{}{
		"Matches": make([]map[string]string, 0),
	}

	dbMatches, _ := db.GetUnfinishedMatches()
	for _, m := range dbMatches {
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

	t.Execute(w, data)
}
