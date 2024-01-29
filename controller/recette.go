package controller

import (
	"groupietracker/back"
	InitTemps "groupietracker/temps"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "index", back.Jeu)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "detail", back.Jeu)
}

func Display(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "display", back.Jeu)
}

func Category(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "category", back.Jeu)
}

func Search(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "search", back.Jeu)
}

func Fav(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "fav", back.Jeu)
}

func Propos(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Url = r.URL.String()
	InitTemps.Temp.ExecuteTemplate(w, "propos", back.Jeu)
}