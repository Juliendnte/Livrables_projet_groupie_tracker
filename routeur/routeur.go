package routeur

import (
	"fmt"
	"groupietracker/back"
	ctrl "groupietracker/controller"
	"groupietracker/temps"
	"net/http"
	"os"
)

func InitServe() {
	http.HandleFunc("/index", ctrl.Index)       //Page Accueil
	http.HandleFunc("/detail", ctrl.Detail)     //Page Detail
	http.HandleFunc("/category", ctrl.Category) //Page Category
	http.HandleFunc("/search", ctrl.Search)     //Page de recherche
	http.HandleFunc("/propos", ctrl.Propos)     //Page a propos
	http.HandleFunc("/favoris", ctrl.Fav)       //Page du compte

	http.HandleFunc("/treatment/favoris", ctrl.InitFav) //Page treatment pour les favoris
	http.HandleFunc("/suppr", ctrl.Suppr)               //Page supprimer d'un favoris
	http.HandleFunc("/img", ctrl.InitImg)                   //Page Img

	http.HandleFunc("/play", ctrl.Play) //Route pour écouter un song
	http.HandleFunc("/url", ctrl.Url)   //Route pour maj le navigator

	http.HandleFunc("/login", ctrl.Login)                           //Page d'affichage des logins
	http.HandleFunc("/login/treatment", ctrl.InitLogin)             //Page de treatment des logins
	http.HandleFunc("/inscription", ctrl.Inscription)               //Page d'affichage des inscriptions
	http.HandleFunc("/inscription/treatment", ctrl.InitInscription) //Page de treatment des inscriptions
	http.HandleFunc("/logout", ctrl.Unlog)                          //Page de déconnexion

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //Page Error 404
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
		temps.Temp.ExecuteTemplate(w, "404", back.Jeu)
	})
	back.Jeu.UtilisateurData.Navigate = back.NewNavigator()

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	fmt.Println("(http://localhost:8080/index) - Server started on port:8080")
	fmt.Println("Si le navigateur ne s'ouvre pas tous seul, tape  http://localhost:8080/index  dans ton navigateur préféré.")
	fmt.Println("Si tu veux arrêter le serveur fait un CTRL + C dans le terminal ")
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Server closed")
}
