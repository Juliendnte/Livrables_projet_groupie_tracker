package main 

import
(
    "groupietracker/temps"
    "groupietracker/routeur"
)

func main(){
    temps.InitTemplate()//Initialisation des templates
    routeur.InitServe()//Initialisation des routes / lancement du serveur 
}