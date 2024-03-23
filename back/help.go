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
	"strconv"
	"strings"
	"time"

	ytb "github.com/kkdai/youtube/v2"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var Token string
var apiKeyYtb = []string{"AIzaSyDf8EHLWiDh1mxWHpzhgEu-7FWK5VZESNg", "AIzaSyBufrRomGBRjjeRmbMW38UO0piC6DG-F5Y", "AIzaSyCk7g4JR_Ea3dywU-v7RiS7S1JNSUFB_Gk", "AIzaSyB3zV564501sl-HN4yWwju1cziIKWVbWv8","AIzaSyBoYLiHCv1UHBgHrhHmFKo8-xoBjm1fjH8","AIzaSyBsKRnoelWNjFGchQ5ynWq8qVz8deQ6RP4"}
var Err error

// Demande a l'api son corps sous un format JSON et le met dans une structure
func RequestApi(apiURL string) ([]byte, ErreurApi) {
	fmt.Println(apiURL)
	//Initialisation du client
	httpClient := http.Client{
		Timeout: time.Second * 10, //Timeout après 10 seconds
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
	fmt.Println("JSON en cours de lecture...")
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
	fmt.Println("JSON en cours de modification...")
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
	fmt.Println("MDP crypté")
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

// Fonction pour savoir si un élément est dans la list
func IsInList(lst []string, s string) bool { // on regarde si une lettre est dans la liste ou pas
	for _, c := range lst {
		if string(c) == s {
			return true
		}
	}
	return false
}

func IsElementPresent(lst1, lst2 []string) bool {
	for _, c1 := range lst1 {
		for _, c2 := range lst2 {
			if c2 == c1 {
				return true
			}
		}
	}
	return false
}

// Fonction pour transformer une liste en string
func TransformSlice(s []string) string { //Met un []string en mot
	var str string
	for _, c := range s {
		str += c + " "
	}
	return str
}

// Fonction pour generer une lettre
func GetRandomLetter() string {
	rand.Seed(time.Now().UnixNano())
	return string('a' + rune(rand.Intn(26)))
}

// Fonction pour generer un chiffre entre 0 et 988
func RandOffset() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(988))
}

// Fonction trie insertion pour le nom
func (arrayToSort *Playlist) InsertionSortPlaylist() {
	fmt.Println("Playlist en cours de trie...")
	for index := 1; index < len(arrayToSort.Playlists.Items); index++ {
		currentItem := arrayToSort.Playlists.Items[index]
		currentLeftIndex := index - 1

		for currentLeftIndex >= 0 && arrayToSort.Playlists.Items[currentLeftIndex].Name > currentItem.Name {
			arrayToSort.Playlists.Items[currentLeftIndex+1] = arrayToSort.Playlists.Items[currentLeftIndex]
			currentLeftIndex -= 1
		}

		arrayToSort.Playlists.Items[currentLeftIndex+1] = currentItem
	}
	fmt.Println("Trie finie")
}

func (arrayToSort *Artists) InsertionSortArtists() {
	fmt.Println("Artists en cours de trie...")
	for index := 1; index < len(arrayToSort.Artists.Items); index++ {
		currentItem := arrayToSort.Artists.Items[index]
		currentLeftIndex := index - 1

		for currentLeftIndex >= 0 && arrayToSort.Artists.Items[currentLeftIndex].Name > currentItem.Name {
			arrayToSort.Artists.Items[currentLeftIndex+1] = arrayToSort.Artists.Items[currentLeftIndex]
			currentLeftIndex -= 1
		}

		arrayToSort.Artists.Items[currentLeftIndex+1] = currentItem
	}
	fmt.Println("Trie finie")
}

func (arrayToSort *Albums) InsertionSortAlbums() {
	fmt.Println("Albums en cours de trie...")
	for index := 1; index < len(arrayToSort.Albums.Items); index++ {
		currentItem := arrayToSort.Albums.Items[index]
		currentLeftIndex := index - 1

		for currentLeftIndex >= 0 && arrayToSort.Albums.Items[currentLeftIndex].Name > currentItem.Name {
			arrayToSort.Albums.Items[currentLeftIndex+1] = arrayToSort.Albums.Items[currentLeftIndex]
			currentLeftIndex -= 1
		}

		arrayToSort.Albums.Items[currentLeftIndex+1] = currentItem
	}
	fmt.Println("Trie finie")
}

func (arrayToSort *Track) InsertionSortTracks() {
	fmt.Println("Track en cours de trie...")
	for index := 1; index < len(arrayToSort.Tracks.Items); index++ {
		currentItem := arrayToSort.Tracks.Items[index]
		currentLeftIndex := index - 1

		for currentLeftIndex >= 0 && arrayToSort.Tracks.Items[currentLeftIndex].Name > currentItem.Name {
			arrayToSort.Tracks.Items[currentLeftIndex+1] = arrayToSort.Tracks.Items[currentLeftIndex]
			currentLeftIndex -= 1
		}

		arrayToSort.Tracks.Items[currentLeftIndex+1] = currentItem
	}
	fmt.Println("Trie finie")
}

func (alb AlbumPrecision) TempsAlbum() string {
	var miliseconds int
	for _, c := range alb.Tracks.Items {
		miliseconds += c.DurationMs
	}
	return Tmps(miliseconds)
}

func (alb PlaylistPrecision) TempsPlaylist() string {
	var miliseconds int
	for _, c := range alb.Tracks.Items {
		miliseconds += c.Track.DurationMs
	}
	return Tmps(miliseconds)
}

// Calcul de miliseconds en heure ou minute
func Tmps(miliseconds int) string {
	var temps Duree

	temps.heure = miliseconds / (1000 * 3600)
	miliseconds %= 1000 * 3600

	temps.min = miliseconds / (1000 * 60)
	miliseconds %= 1000 * 60

	temps.sec = miliseconds / 1000

	return formatDuree(temps)
}

// Affichage du temps comme sur spotify
func formatDuree(temps Duree) string {
	if temps.heure > 0 {
		return strconv.Itoa(temps.heure) + " heure " + strconv.Itoa(temps.min) + " min "
	}
	if len(strconv.Itoa(temps.sec)) == 1 {
		return strconv.Itoa(temps.min) + ":0" + strconv.Itoa(temps.sec)
	}
	return strconv.Itoa(temps.min) + ":" + strconv.Itoa(temps.sec)
}

// Affichage des abonnement
func FormatAbo(number int) string {
	numStr := strconv.Itoa(number)
	if len(numStr) <= 3 {
		return numStr
	}

	var result string
	for i := 0; i < len(numStr); i++ {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result += " "
		}
		result += string(numStr[i])
	}

	return result
}

// Téléchargement d'une vidéo ytb
func Download(videoID string) (string, error) {
	fmt.Println("La vidéo " + videoID + " est en cours de chargement")
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
	fmt.Println("Fin du chargement")
	return videoID, nil
}

// Retourne le nombre de like d'une vidéo sur ytb
func Like(searchTerm string) string {
	if len(apiKeyYtb) <= 0 {
		fmt.Println("Erreur lors de la recherche de vidéos : googleapi: Error 403: The request cannot be completed because you have exceeded your <a href='/youtube/v3/getting-started#quota'>quota</a>., quotaExceeded")
		return "0"
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKeyYtb[0]},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Printf("Erreur lors de la création du client YouTube : %v\n", err)
		return "0"
	}

	searchResponse, err := service.Search.List([]string{"id", "snippet"}).Q(searchTerm).Type("video").MaxResults(1).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Error 403") {
			apiKeyYtb = apiKeyYtb[1:]
			fmt.Println("Changement de compte")
			return Like(searchTerm)
		}
		fmt.Printf("Erreur lors de la recherche de vidéos : %v\n", err)
		return "0"
	}

	if len(searchResponse.Items) == 0 {
		fmt.Printf("Aucune vidéo trouvée pour le terme de recherche : %s\n", searchTerm)
		return "0"
	}

	if videoStatsResponse, err := service.Videos.List([]string{"statistics"}).Id(searchResponse.Items[0].Id.VideoId).Do(); err != nil {
		fmt.Printf("Erreur lors de la récupération des statistiques de la vidéo : %v\n", err)
		return "0"
	} else {
		return FormatAbo(int(videoStatsResponse.Items[0].Statistics.LikeCount))
	}
}

// Retourne l'id d'une vidéo sur ytb
func IdYtb(search string) string {
	if len(apiKeyYtb) <= 0 {
		fmt.Println("Erreur lors de la recherche de vidéos : googleapi: Error 403: The request cannot be completed because you have exceeded your <a href='/youtube/v3/getting-started#quota'>quota</a>., quotaExceeded")
		return "0"
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKeyYtb[0]},
	}

	service, err := youtube.New(client)
	if err != nil {
		if strings.Contains(err.Error(), "Error 403") {
			apiKeyYtb = apiKeyYtb[1:]
			fmt.Println("Changement de compte")
			return Like(search)
		}
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
