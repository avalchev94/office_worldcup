package main

import (
	"github.com/avalchev94/office_worldcup/database"
	"net/http"
	"path/filepath"
	"text/template"
)

type predictionData struct {
	Host      string
	Guest     string
	Predicted string
	Result    string
	Score     int32
}

func predictionsHandler(w http.ResponseWriter, r *http.Request) {
	if user, err := authenticated(w, r); err == nil {
		db, _ := database.New()
		defer db.Close()

		dbPredictions, _ := db.GetPredictionsByUser(user.ID)
		predictions := make([]predictionData, len(dbPredictions))
		for i, p := range dbPredictions {
			m, _ := db.GetMatch(p.Match)
			host, _ := db.GetTeamByID(m.Host)
			guest, _ := db.GetTeamByID(m.Guest)

			predictions[i] = predictionData{
				Host:      host.Name,
				Guest:     guest.Name,
				Predicted: p.Predicted,
				Result:    m.Result,
				Score:     p.Score,
			}
		}

		fullpath := filepath.Join(templateFolder, "predictions.html")
		t := template.Must(template.ParseFiles(fullpath))
		t.Execute(w, map[string]interface{}{
			"Predictions": predictions,
		})
	}
}
