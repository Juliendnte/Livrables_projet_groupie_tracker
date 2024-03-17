package routeur

import (
	"fmt"
	ctrl "groupietracker/controller"
	"net/http"
	"os"
)

func InitServe() {

	http.HandleFunc("/login", ctrl.Login)
	http.HandleFunc("/login/treatment", ctrl.InitLogin)
	http.HandleFunc("/inscription", ctrl.Inscription)
	http.HandleFunc("/inscription/treatment", ctrl.InitInscription)
	http.HandleFunc("/logout", ctrl.Unlog)

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	fmt.Println("(http://localhost:8080/index) - Server started on port:8080")
	fmt.Println("Si le navigateur ne c'est pas ouvrir tous seul, va y tous seul et tape  http://localhost:8080/index  dans ton navigateur préféré.")
	fmt.Println("Si tu veux arrêter le server fait un CTRL + C dans le terminal ")
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Server closed")
}
