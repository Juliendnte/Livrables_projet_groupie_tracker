package routeur

import (
	"fmt"
	ctrl "groupietracker/controller"
	"net/http"
	"os"
)

func InitServe() {

	http.HandleFunc("/index", ctrl.Index)
	http.HandleFunc("/display", ctrl.Display)
	http.HandleFunc("/detail", ctrl.Detail)
	http.HandleFunc("/category", ctrl.Category)
	http.HandleFunc("/search", ctrl.Search)
	http.HandleFunc("/favoris", ctrl.Fav)
	http.HandleFunc("/propos", ctrl.Propos)

	http.HandleFunc("/inscription", ctrl.Login)
	http.HandleFunc("/login", ctrl.Inscription)

	http.HandleFunc("/", ctrl.HandleError)

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	fmt.Println("(http://localhost:8080/) - Server started on port:8080")
	fmt.Println("Si le navigateur ne c'est pas ouvrir tous seul, va y tous seul et tape  http://localhost:8080/index  dans ton navigateur préféré.")
	fmt.Println("Si tu veux arrêter le server fait un CTRL + C dans le terminal ")
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Server closed")
}
