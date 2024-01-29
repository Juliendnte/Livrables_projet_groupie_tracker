package back

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Demande a l'api son corps sous un format JSON et le met dans une structure
func RequestApi(apiURL string) {
	//Initialisation du client
	httpClient := http.Client{
		Timeout: time.Second * 2, //Timeout après 2 seconds
	}
	//Création de la reguête HTTP avec un GET vers l'api
	req, errReq := http.NewRequest("GET", apiURL, nil)
	if errReq != nil {
		fmt.Println("Error creating request :", errReq.Error())
		os.Exit(1)
	}
	//Ajout du token dans l'header pour avoir l'autorisation d'émettre des requettes sur l'api de spotify
	req.Header.Set("User_Agent", "Question pour un clanpin")

	//Exécution de la requête HTTP vers l'api
	resp, errRes := httpClient.Do(req)
	if resp.Body != nil {
		defer resp.Body.Close()
	} else {
		fmt.Println("Error creating response :", errRes.Error())
		os.Exit(2)
	}

	//Lecture du corps de la requête HTTP
	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println("Error reading response body :", errBody.Error())
		os.Exit(3)
	}

	json.Unmarshal(body, &lstNourriture)
}

// Fonction pour récupèrer le mot de passe crypté
func MdpCrypt(Mdp string) string {
	jsonFile, err := os.ReadFile("JSON/login.json") //Récupére les données du JSON
	if err != nil {
		fmt.Println("Error reading", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(jsonFile, &lstNourriture) //Met dans ma struct
	if err != nil {
		fmt.Println("Error encodage ", err.Error())
		os.Exit(1)
	}

	hasher := sha256.New()
	hasher.Write([]byte(Mdp))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword // mdp crypter
}
