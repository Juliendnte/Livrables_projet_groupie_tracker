package controller

import (
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
	back.Jeu.UtilisateurData.Fav = false
	back.Jeu.Utilisateur = back.Client{}
	InitTemp.Temp.ExecuteTemplate(w, "unlog", back.Jeu)
}

// Fonction pour se connecter
func Login(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		}
		return
	}
	InitTemp.Temp.ExecuteTemplate(w, "Login", back.Jeu)
}

// Fonction pour s'inscrire
func Inscription(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		}
		return
	}
	InitTemp.Temp.ExecuteTemplate(w, "inscription", back.Jeu)
}

// Fonction treatment pour se connecter
func InitLogin(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		}
		return
	}

	back.User.Name = r.FormValue("Nom")
	back.User.Mdp = back.MdpCrypt(r.FormValue("Mdp")) //Récupére les données de l'utilisateur

	for _, c := range back.LstUser {
		if back.User.Name == c.Name {
			if back.User.Mdp == c.Mdp {
				back.Jeu.UtilisateurData.Connect = true
				if len(c.FavorisAlbum) == 0 && len(c.FavorisArtist) == 0 && len(c.FavorisPlaylist) == 0 && len(c.FavorisTrack) == 0 {
					back.Jeu.UtilisateurData.Fav = false
				} else {
					back.Jeu.UtilisateurData.Fav = true
				}
				back.Jeu.Utilisateur = c
				if !Header(w, true) {
					return
				}
				if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
					http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
				} else {
					http.Redirect(w, r, "/index", http.StatusMovedPermanently)
				}
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
		if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		}
		return
	}

	if back.LstUser, err = back.ReadJSON(); err != nil {
		return
	}

	back.User.Name = r.FormValue("Nom")
	back.User.Mdp = back.MdpCrypt(r.FormValue("Mdp")) //Récupére les données de l'utilisateur
	back.User.Id = back.GenerateID()
	back.User.Img = "vert"
	back.User.Letter = string(back.User.Name[0])

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
	back.EditJSON(back.LstUser)
	back.Jeu.UtilisateurData.Connect = true
	back.Jeu.Utilisateur = back.User
	if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
}

// Fonction pour engregistrer une image et retourne son nom
func InitImg(w http.ResponseWriter, r *http.Request) {
	if !back.Jeu.UtilisateurData.Connect {
		if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		}
		return
	}
	//Prend les données ne dépassant cette taille (pout l'image)
	if r.ParseMultipartForm(10<<20) != nil {
		fmt.Println("Error r.ParseMultipartForm")
		http.Redirect(w, r, "/favoris", http.StatusMovedPermanently)
		return
	}

	file, handler, errFile := r.FormFile("Img") //Récupère le fichier image
	if errFile != nil {
		fmt.Printf(errFile.Error(), "\n")
		http.Redirect(w, r, "/favoris", http.StatusMovedPermanently)
		return
	}
	defer file.Close()
	f, _ := os.Create("./assets/img/" + handler.Filename) //Chemin où mettre le fichier
	defer f.Close()
	io.Copy(f, file) //Met l'image au chemin donnée
	back.Jeu.Utilisateur.Img = handler.Filename
	if back.LstUser, err = back.ReadJSON(); err != nil {
		return
	}
	for i, c := range back.LstUser {
		if c.Id == back.Jeu.Utilisateur.Id {
			back.LstUser[i].Img = back.Jeu.Utilisateur.Img
		}
	}
	back.EditJSON(back.LstUser)
	http.Redirect(w, r, "/favoris", http.StatusMovedPermanently)
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
	if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
}
