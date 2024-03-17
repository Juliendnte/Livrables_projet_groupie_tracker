package controller

import (
	"encoding/json"
	"fmt"
	back "groupietracker/back"
	InitTemp "groupietracker/temps"
	"io"
	"net/http"
	"os"
	"strings"
)

var err error

// Fonction pour se déconnecter
func Unlog(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Connect = false
	back.Jeu.Utilisateur = back.Client{}
	InitTemp.Temp.ExecuteTemplate(w, "unlog", back.Jeu)
}

// Fonction pour se connecter
func Login(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		return
	}
	//back.Jeu.UtilisateurData = back.UserData
	InitTemp.Temp.ExecuteTemplate(w, "Login", back.Jeu)
}

// Fonction pour s'inscrire
func Inscription(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		return
	}
	//back.Jeu.UtilisateurData = back.UserData
	InitTemp.Temp.ExecuteTemplate(w, "inscription", back.Jeu)
}

// Fonction treatment pour se connecter
func InitLogin(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		return
	}

	back.User.Name = r.FormValue("Nom")
	back.User.Mdp = back.MdpCrypt(r.FormValue("Mdp")) //Récupére les données de l'utilisateur

	for _, c := range back.LstUser {
		if back.User.Name == c.Name {
			if back.User.Mdp == c.Mdp {
				back.Jeu.UtilisateurData.Connect = true
				back.Jeu.Utilisateur = c
				http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
				return
			}
		}
	}
	//Sinon reste sur la page login
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// Fonction treatment pour se connecter
func InitInscription(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		return
	}

	back.User.Name = r.FormValue("Nom")
	back.User.Mdp = back.MdpCrypt(r.FormValue("Mdp")) //Récupére les données de l'utilisateur
	back.User.Id = back.GenerateID()
	back.User.Img = InitImg(w, r)
	if back.User.Img == "vert" {
		back.User.Letter = string(back.User.Name[0])
	}

	for _, c := range back.LstUser {
		if back.User.Name == c.Name {
			if back.User.Mdp == c.Mdp {
				fmt.Println(back.User.Name, c.Name)
				fmt.Println(back.User.Mdp, c.Mdp)
				http.Error(w, "Username already exists", http.StatusConflict)
				return
			}
		}
	}

	back.LstUser = append(back.LstUser, back.User) //Ajoute l'utilisateur

	if back.Body, err = json.Marshal(back.LstUser); err != nil {
		fmt.Println("-----------------Error encodage InitInscription-----------------", err.Error())
		return
	}

	// Écrire le JSON modifié dans le fichier
	if os.WriteFile("JSON/login.json", back.Body, 0644) != nil {
		fmt.Println("-----------------Erreur lors de l'écriture du fichier JSON modifié-----------------")
		return
	}

	back.Jeu.UtilisateurData.Connect = true
	back.Jeu.Utilisateur = back.User
	http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
}

// Fonction pour engregistrer une image et retourne son nom
func InitImg(w http.ResponseWriter, r *http.Request) string {
	if back.Jeu.UtilisateurData.Connect {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		return ""
	}
	//Prend les données ne dépassant cette taille (pout l'image)
	if r.ParseMultipartForm(10<<20) != nil {
		return "vert"
	}

	file, handler, errFile := r.FormFile("Img") //Récupère le fichier image
	if errFile != nil {
		return "vert"
	}
	defer file.Close()
	f, _ := os.Create("./assets/img/" + handler.Filename) //Chemin où mettre le fichier
	defer f.Close()
	io.Copy(f, file) //Met l'image au chemin donnée

	return handler.Filename
}

func Url(w http.ResponseWriter, r *http.Request) {
	var link string
	if r.URL.Query().Get("url") == "before" {
		link = back.Jeu.UtilisateurData.Navigate.GoBack()
	} else {
		link = back.Jeu.UtilisateurData.Navigate.GoForward()
	}
	if link != "No more" {
		if strings.Contains(link, "?") {
			link += "&url=" + r.URL.Query().Get("url")
		} else {
			link += "?url=" + r.URL.Query().Get("url")
		}
		http.Redirect(w, r, link, http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
}
