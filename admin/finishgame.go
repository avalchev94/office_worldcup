package admin

import (
	"github.com/avalchev94/office_worldcup/database"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
)

func calculateScore(predicted, expected string) int32 {
	getSign := func(result string) rune {
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

	if predicted == expected {
		return 4
	} else if getSign(predicted) == getSign(expected) {
		return 2
	} else {
		return -1
	}
}

func finishGameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	regex, _ := regexp.Compile("[0-9]+:[0-9]+")
	result := r.Form["result"][0]

	if regex.Match([]byte(result)) {
		db, _ := database.New()
		defer db.Close()

		matchID := database.ObjectIdHex(mux.Vars(r)["id"])
		dbPredictions, err := db.GetPredictionsByMatch(matchID)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		for _, p := range dbPredictions {
			p.Score = calculateScore(p.Predicted, result)
			db.UpdatePrediction(p)
		}

		match, _ := db.GetMatch(matchID)
		match.Result = result
		db.UpdateMatch(match)

		http.Redirect(w, r, "/admin", http.StatusPermanentRedirect)
	} else {
		w.Write([]byte("incorrect result"))
	}
}
