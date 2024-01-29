package controller

import (
	InitTemp "groupietracker/temps"
	back "groupietracker/back"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var err error

// Fonction pour se déconnecter
func Unlog(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Connect = false
	back.Jeu.Utilisateur = back.Client{"", ""}
	http.Redirect(w, r, back.UserData.Url, http.StatusMovedPermanently)
}

// Fonction pour se connecter
func Login(w http.ResponseWriter, r *http.Request) {
	//back.Jeu.UtilisateurData = back.UserData
	InitTemp.Temp.ExecuteTemplate(w, "Login", back.Jeu)
}

// Fonction pour s'inscrire
func Inscription(w http.ResponseWriter, r *http.Request) {
	//back.Jeu.UtilisateurData = back.UserData
	InitTemp.Temp.ExecuteTemplate(w, "inscription", back.Jeu)
}

// Fonction treatment pour se connecter
func InitLogin(w http.ResponseWriter, r *http.Request) {
	back.User.Name = r.FormValue("Nom")
	back.User.Mdp = back.MdpCrypt(r.FormValue("Mdp")) //Récupére les données de l'utilisateur

	for _, c := range back.LstUser {
		if back.User.Name == c.Name {
			if back.User.Mdp == c.Mdp {
				back.UserData.Connect = true //Le connecte
				back.Jeu.UtilisateurData = back.UserData
				back.Jeu.Utilisateur = c //Lui met ses droits
				http.Redirect(w, r, back.Jeu.UtilisateurData.Url, http.StatusMovedPermanently)
				return
			}
		}
	}
	//Sinon reste sur la page login
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// Fonction treatment pour se connecter
func InitInscription(w http.ResponseWriter, r *http.Request) {
	back.User.Name = r.FormValue("Nom")
	back.User.Mdp = back.MdpCrypt(r.FormValue("Mdp")) //Récupére les données de l'utilisateur

	for _, c := range back.LstUser {
		if back.User.Name == c.Name {
			if back.User.Mdp == c.Mdp {
				http.Error(w, "Username already exists", http.StatusConflict)
				return
			}
		}
	}

	back.LstUser = append(back.LstUser, back.User) //Ajoute l'utilisateur

	modifiedJSON, errMarshal := json.Marshal(back.LstUser) //Met la struct en JSON file
	if errMarshal != nil {
		fmt.Println("Error encodage ", errMarshal.Error())
		return
	}

	// Écrire le JSON modifié dans le fichier
	err = os.WriteFile("JSON/login.json", modifiedJSON, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier JSON modifié:", err)
		return
	}

	back.UserData.Connect = true
	back.Jeu.UtilisateurData = back.UserData
	back.Jeu.Utilisateur = back.User
	http.Redirect(w, r, back.Jeu.UtilisateurData.Url, http.StatusMovedPermanently)
}

// Fonction pour la page error 404
func HandleError(w http.ResponseWriter, r *http.Request) {
	InitTemp.Temp.ExecuteTemplate(w, "error", back.Jeu)
}
