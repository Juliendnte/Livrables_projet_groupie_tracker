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

	http.HandleFunc("/index", ctrl.Index)
	// http.HandleFunc("/display", ctrl.Display)
	http.HandleFunc("/detail", ctrl.Detail)
	http.HandleFunc("/category", ctrl.Category)
	http.HandleFunc("/search", ctrl.Search)
	http.HandleFunc("/treatment/favoris", ctrl.InitFav)
	http.HandleFunc("/suppr", ctrl.Suppr)
	http.HandleFunc("/favoris", ctrl.Fav)
	http.HandleFunc("/propos", ctrl.Propos)
	http.HandleFunc("/play", ctrl.Play)
	http.HandleFunc("/url", ctrl.Url)

	http.HandleFunc("/login", ctrl.Login)
	http.HandleFunc("/login/treatment", ctrl.InitLogin)
	http.HandleFunc("/inscription", ctrl.Inscription)
	http.HandleFunc("/inscription/treatment", ctrl.InitInscription)
	http.HandleFunc("/logout", ctrl.Unlog)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
		temps.Temp.ExecuteTemplate(w, "404", back.Jeu)
	})

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	fmt.Println("(http://localhost:8080/index) - Server started on port:8080")
	fmt.Println("Si le navigateur ne c'est pas ouvrir tous seul, va y tous seul et tape  http://localhost:8080/index  dans ton navigateur préféré.")
	fmt.Println("Si tu veux arrêter le server fait un CTRL + C dans le terminal ")
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Server closed")
}
