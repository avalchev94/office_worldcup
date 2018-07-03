package admin

import (
	"github.com/avalchev94/office_worldcup/database"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
)

func getResultSign(result string) rune {
	l := strings.Split(result, ":")
	switch {
	case l[0] == l[1]:
		return 'X'
	case l[0] > l[1]:
		return '1'
	default:
		return '2'
	}
}

func calculateGroupScore(predicted, expected string) int32 {
	if predicted == expected {
		return 4
	} else if getResultSign(predicted) == getResultSign(expected) {
		return 2
	} else {
		return -1
	}
}

func calculateKnockoutScore(predicted, expected, winner string) int32 {
	prediction := strings.Split(predicted, ";")

	score := calculateGroupScore(prediction[0], expected)
	if score > 0 && getResultSign(expected) == 'X' {
		score++
		if prediction[1] == winner {
			score++
		} else {
			score--
		}
	}

	return score
}

func finishGameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	regex, _ := regexp.Compile("[0-9]+:[0-9]+")
	result := r.FormValue("result")

	if regex.Match([]byte(result)) {
		db, _ := database.New()
		defer db.Close()

		matchID := database.ObjectIdHex(mux.Vars(r)["id"])
		match, _ := db.GetMatch(matchID)
		dbPredictions, _ := db.GetPredictionsByMatch(matchID)

		for _, p := range dbPredictions {
			switch match.Stage {
			case database.GroupStage:
				p.Score = calculateGroupScore(p.Predicted, result)
			case database.KnockoutStage:
				winner := r.FormValue("winner")
				p.Score = calculateKnockoutScore(p.Predicted, result, winner)
			}
			db.UpdatePrediction(p)
		}

		switch match.Stage {
		case database.GroupStage:
			match.Result = result
		case database.KnockoutStage:
			match.Result = result + ";" + r.FormValue("winner")
		}
		db.UpdateMatch(match)

		http.Redirect(w, r, "/admin", http.StatusPermanentRedirect)
	} else {
		w.Write([]byte("incorrect result"))
	}
}
