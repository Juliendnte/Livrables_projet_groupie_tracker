package back

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var Token string
var Err error

// Demande a l'api son corps sous un format JSON et le met dans une structure
func RequestApi(apiURL string) ([]byte, ErreurApi) {
	//fmt.Println(apiURL)
	//Initialisation du client
	httpClient := http.Client{
		Timeout: time.Second * 4, //Timeout après 4 seconds
	}
	//Création de la reguête HTTP avec un GET vers l'api
	req, errReq := http.NewRequest("GET", apiURL, nil)
	if errReq != nil {
		fmt.Println("-----------------Error creating request :-----------------", errReq.Error())
		return nil, Fail
	}
	//Ajout du token dans l'header pour avoir l'autorisation d'émettre des requettes sur l'api de spoonacular
	req.Header.Set("Authorization", "Bearer "+Token)

	//Exécution de la requête HTTP vers l'api
	resp, errRes := httpClient.Do(req)
	if errRes != nil {
		fmt.Println("-----------------Error creating response :-----------------", errRes.Error())
		Fail.Error.Status = 503
		Fail.Error.Message = "No Service Unavailable"
		return nil, Fail
	} else {
		defer resp.Body.Close()
	}

	//Lecture du corps de la requête HTTP
	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println("-----------------Error reading response body :-----------------", errBody.Error())
		return nil, Fail
	}

	if resp.StatusCode != 200 {
		json.Unmarshal(body, &Fail)
		if Fail.Error.Message == "Only valid bearer authentication supported" || Fail.Error.Status == 401 {
			fmt.Println(Fail)
			ReloadApi()
			return RequestApi(apiURL)
		}
	} else {
		Fail.Error.Status = resp.StatusCode
		Fail.Error.Message = resp.Status
	}
	return body, Fail
}

// Recharge le token s'il n'est plus bon d'accès (token usage de 1H)
func ReloadApi() {
	//URL de l'api de Spotify pour avoir le token
	urlToken := "https://accounts.spotify.com/api/token"
	const clientId = "b2717725a58e4d6faccd1ee5fd5bd55b" //Pas touche les miens
	const clientSecret = "67c6ede35f3846bc95c2093a4d6e232c"

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	//Création du corps de ma requête pour avoir les accès au token
	BodyReq := bytes.NewBufferString("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientSecret)

	///Création de la reguête HTTP vers l'api en POST avec l'url du token et le corps de ma requête
	req, errReq := http.NewRequest("POST", urlToken, BodyReq)
	if errReq != nil {
		fmt.Println("-----------------Error creating request :-----------------", errReq.Error())
		return
	}
	//Metadonné nécessaire dans le header pour une requête POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, errRes := httpClient.Do(req)
	if resp.Body != nil {
		defer resp.Body.Close()
	} else {
		fmt.Println("-----------------Error creating response :-----------------", errRes.Error())
		return
	}
	//Décodage du JSON dans une map
	var repMap map[string]interface{}

	decoder := json.NewDecoder(resp.Body)
	if errJSON := decoder.Decode(&repMap); errJSON != nil {
		fmt.Println("-----------------Error reading JSON : -----------------", errJSON.Error())
		return
	}
	//Le nouveau token est mit dans la variable globale
	Token = repMap["access_token"].(string)
	fmt.Println("Bearer ", Token)
}
