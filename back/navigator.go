package back

import "fmt"

type Navigator struct {
	History []string
	Index   int
}

//Initialisation du système de navigation
func NewNavigator() *Navigator {
	fmt.Println("Le système de navigateur est initialisé...")
	return &Navigator{
		History: make([]string, 0),
		Index:   -1,
	}
}

//Visiter une page
func (nav *Navigator) VisitPage(page string) {
	if nav.Index < len(nav.History)-1 {
		nav.History = nav.History[:nav.Index+1]
	}
	nav.History = append(nav.History, page)
	fmt.Println("L'utilisateur a visité ", page)
	nav.Index++
}

//Retour sur une ancienne page 
func (nav *Navigator) GoBack() string {
	if nav.Index <= 0 {
		return "No more"
	}
	nav.Index--
	fmt.Println("L'utilisateur a fait un retour en arrière vers ", nav.History[nav.Index])
	return nav.History[nav.Index]
}

//Aller sur une page vers l'avant
func (nav *Navigator) GoForward() string {
	if nav.Index >= len(nav.History)-1 {
		return "No more"
	}
	nav.Index++
	fmt.Println("L'utilisateur est allé en avant vers ", nav.History[nav.Index])
	return nav.History[nav.Index]
}

//Revenire sur notre page actuel
func (nav *Navigator) GoAcutal() string {
	fmt.Println("L'utilisateur est actuellement sur la page ", nav.History[len(nav.History)-1])
	return nav.History[len(nav.History)-1]
}
