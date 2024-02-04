package main 

import
(
    "groupietracker/temps"
    "groupietracker/routeur"
)

func main(){
    temps.InitTemplate()
    routeur.InitServe()
}