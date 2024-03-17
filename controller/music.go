package controller

import (
	"encoding/json"
	"fmt"
	"groupietracker/back"
	InitTemps "groupietracker/temps"
	"net/http"
	"strconv"
)

func Header(w http.ResponseWriter) bool {
	if back.Jeu.Header.Albums.Href == "" {
		if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/browse/new-releases?limit=5"); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return false
		}
		json.Unmarshal(back.Body, &back.Jeu.Header)
	}
	return true
}
