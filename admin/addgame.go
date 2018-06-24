package admin

import (
	"github.com/avalchev94/office_worldcup/database"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func parseTime(dateStr, timeStr string) time.Time {
	dateList := make([]int, 3)
	for i, d := range strings.Split(dateStr, "-") {
		dateList[i], _ = strconv.Atoi(d)
	}
	timeList := make([]int, 2)
	for i, t := range strings.Split(timeStr, ":") {
		timeList[i], _ = strconv.Atoi(t)
	}

	return time.Date(dateList[0], time.Month(dateList[1]), dateList[2], timeList[0], timeList[1], 0, 0, time.Local).UTC()
}

func addGameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fullpath := filepath.Join(templateFolder, "add_game.html")
		t := template.Must(template.ParseFiles(fullpath))

		db, _ := database.New()
		defer db.Close()

		dbTeams, _ := db.GetTeams()
		data := map[string]interface{}{
			"Teams": make([]map[string]string, len(dbTeams)),
		}
		for i, t := range dbTeams {
			teams := data["Teams"].([]map[string]string)
			teams[i] = map[string]string{
				"ID":   t.ID.Hex(),
				"Name": t.Name,
			}
		}

		t.Execute(w, data)
	case http.MethodPost:
		r.ParseForm()

		db, _ := database.New()
		defer db.Close()

		match := database.Match{
			Host:   database.ObjectIdHex(r.FormValue("host")),
			Guest:  database.ObjectIdHex(r.FormValue("guest")),
			Result: "",
			Date:   parseTime(r.FormValue("date"), r.FormValue("time")),
		}

		db.AddMatch(match)
	}
}
