package back

type Site struct{
	Nourriture Food
	Utilisateur Client
	UtilisateurData ClientData

}

type Food struct {
	Vegetarian               bool    `json:"vegetarian"`
	Vegan                    bool    `json:"vegan"`
	GlutenFree               bool    `json:"glutenFree"`
	DairyFree                bool    `json:"dairyFree"`
	VeryHealthy              bool    `json:"veryHealthy"`
	Cheap                    bool    `json:"cheap"`
	VeryPopular              bool    `json:"veryPopular"`
	PreparationMinutes       int     `json:"preparationMinutes"`
	CookingMinutes           int     `json:"cookingMinutes"`
	HealthScore              int     `json:"healthScore"`
	PricePerServing          float64 `json:"pricePerServing"`
	ExtendedIngredients      []struct {
		ID           int           `json:"id"`
		Aisle        string        `json:"aisle"`
		Image        string        `json:"image"`
		Name         string        `json:"name"`
		Measures     struct {
			Metric struct {
				Amount    float64 `json:"amount"`
				UnitShort string  `json:"unitShort"`
			} `json:"metric"`
		} `json:"measures"`
	} `json:"extendedIngredients"`
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	ReadyInMinutes int           `json:"readyInMinutes"`
	SourceURL      string        `json:"sourceUrl"`
	Image          string        `json:"image"`
	Summary        string        `json:"summary"`
	Cuisines       []interface{} `json:"cuisines"`
	DishTypes      []string      `json:"dishTypes"`
	Diets          []string      `json:"diets"`
	Occasions      []interface{} `json:"occasions"`
	WinePairing    struct {
		PairedWines    []string `json:"pairedWines"`
		PairingText    string   `json:"pairingText"`
	} `json:"winePairing"`
	Instructions         string `json:"instructions"`
}
type ClientData struct {
	Url string 
	Connect bool
}

type Client struct {
	Name  string `json:"name"`
	Mdp   string `json:"mdp"`
}

var lstNourriture Food
var LstUser []Client
var User Client
var UserData ClientData
var Jeu Site