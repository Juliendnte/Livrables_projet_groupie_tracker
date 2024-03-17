package back

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	ytb "github.com/kkdai/youtube/v2"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var Token string
var apiKeyYtb = "AIzaSyDf8EHLWiDh1mxWHpzhgEu-7FWK5VZESNg"
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

// Fonction pour mettre le JSON dans une struct
func ReadJSON() ([]Client, error) {
	jsonFile, err := os.ReadFile("JSON/login.json")
	if err != nil {
		fmt.Println("-----------------Error reading-----------------", err.Error())
	}

	var jsonData []Client
	err = json.Unmarshal(jsonFile, &jsonData)
	return jsonData, err
}

// Fonction pour modifié le JSON
func EditJSON(ModifiedClient []Client) {

	modifiedJSON, errMarshal := json.Marshal(ModifiedClient)
	if errMarshal != nil {
		fmt.Println("-----------------Error encodage -----------------", errMarshal.Error())
		return
	}

	// Écrire le JSON modifié dans le fichier
	if err := os.WriteFile("JSON/login.json", modifiedJSON, 0644); err != nil {
		fmt.Println("-----------------Erreur lors de l'écriture du fichier JSON modifié:-----------------", err)
	}
}

// Fonction pour récupèrer le mot de passe crypté
func MdpCrypt(Mdp string) string {
	jsonFile, err := os.ReadFile("JSON/login.json") //Récupére les données du JSON
	if err != nil {
		fmt.Println("-----------------Error reading MdpCrypt-----------------", err.Error())
		return err.Error()
	}

	if err = json.Unmarshal(jsonFile, &LstUser); err != nil {
		fmt.Println("-----------------Error encodage MdpCrypt-----------------", err.Error())
		return err.Error()
	}

	hasher := sha256.New()
	hasher.Write([]byte(Mdp))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword // mdp crypter
}

// Fonction pour savoir si l'id existe déjà
func IdAlreadyExists(nb int) bool {
	for i := 0; i < len(LstUser); i++ {
		if LstUser[i].Id == nb {
			return true
		}
	}
	return false
}

// Fonction pour générer un Id disponible
func GenerateID() int {
	var Id int = rand.Intn(100)
	if IdAlreadyExists(Id) {
		return GenerateID()
	}
	return Id
}

// Téléchargement d'une vidéo ytb
func Download(videoID string) (string, error) {
	client := ytb.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		fmt.Println("Error getting video ", err)
		return "", err
	}

	stream, _, err := client.GetStream(video, &video.Formats.WithAudioChannels()[0])
	if err != nil {
		fmt.Println("Error getting stream ", err)
		return "", err
	}
	defer stream.Close()

	file, err := os.Create("./assets/vid/" + videoID + ".mp3")
	if err != nil {
		fmt.Println("Error creating .mp4 at .assets/vid/ ", err)
		return "", err
	}
	defer file.Close()

	if _, err = io.Copy(file, stream); err != nil {
		fmt.Println("Error putting stream at file ", err)
		return "", err
	}

	return videoID, nil
}

// Retourne le nombre de like d'une vidéo sur ytb
func Like(searchTerm string) int {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKeyYtb},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Printf("Erreur lors de la création du client YouTube : %v\n", err)
		return 0
	}

	searchResponse, err := service.Search.List([]string{"id", "snippet"}).Q(searchTerm).Type("video").MaxResults(1).Do()
	if err != nil {
		fmt.Printf("Erreur lors de la recherche de vidéos : %v\n", err)
		return 0
	}

	if len(searchResponse.Items) == 0 {
		fmt.Printf("Aucune vidéo trouvée pour le terme de recherche : %s\n", searchTerm)
		return 0
	}

	if videoStatsResponse, err := service.Videos.List([]string{"statistics"}).Id(searchResponse.Items[0].Id.VideoId).Do(); err != nil {
		fmt.Printf("Erreur lors de la récupération des statistiques de la vidéo : %v\n", err)
		return 0
	} else {
		return int(videoStatsResponse.Items[0].Statistics.LikeCount)
	}
}

// Retourne l'id d'une vidéo sur ytb
func IdYtb(search string) string {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKeyYtb},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Printf("Erreur lors de la création du client YouTube : %v", err)
		return "dQw4w9WgXcQ" //Rick Roll
	}

	searchResponse, err := service.Search.List([]string{"id", "snippet"}).Q(search).Type("video").MaxResults(1).Do()
	if err != nil {
		fmt.Printf("Erreur lors de la recherche de vidéos : %v", err)
		return "dQw4w9WgXcQ"
	}

	if len(searchResponse.Items) == 0 {
		fmt.Printf("Aucune vidéo trouvée pour le terme de recherche : %s", search)
		return "dQw4w9WgXcQ"
	}

	return searchResponse.Items[0].Id.VideoId
}
