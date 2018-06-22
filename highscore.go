package main

import (
	"github.com/avalchev94/office_worldcup/database"
	"net/http"
	"path/filepath"
	"sort"
	"text/template"
)

type userData struct {
	Username    string
	Fullname    string
	Score       int
	Predictions int
}

func highscoreHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := authenticated(w, r); err == nil {
		db, _ := database.New()
		defer db.Close()

		users := make([]userData, 0)
		dbUsers, _ := db.GetUsers()
		for _, u := range dbUsers {
			current := userData{
				Username:    u.Username,
				Fullname:    u.Fullname,
				Score:       0,
				Predictions: 0,
			}

			if dbPredictions, err := db.GetPredictionsByUser(u.ID); err == nil {
				current.Predictions = len(dbPredictions)

				for _, p := range dbPredictions {
					current.Score += int(p.Score)
				}
			}

			users = append(users, current)
		}

		sort.Slice(users, func(i, j int) bool {
			if users[i].Score != users[j].Score {
				return users[i].Score > users[j].Score
			}

			if users[i].Predictions != users[j].Predictions {
				return users[i].Predictions > users[j].Predictions
			}

			return users[i].Username > users[j].Username
		})

		fullpath := filepath.Join(templateFolder, "highscore.html")
		t := template.Must(template.ParseFiles(fullpath))
		t.Execute(w, map[string]interface{}{
			"Users": users,
		})
	}
}
