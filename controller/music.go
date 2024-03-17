package controller

import (
	"encoding/json"
	"fmt"
	back "groupietracker/back"
	InitTemps "groupietracker/temps"
	"net/http"
	"strconv"
)

var limit = 12

func Header(w http.ResponseWriter) bool {
	if back.Jeu.Header.Albums.Href == "" {
		if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/browse/new-releases?limit=5"); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return false
		}
		json.Unmarshal(back.Body, &back.Jeu.Header)
	}
	return true
}

func Index(w http.ResponseWriter, r *http.Request) {
	back.Jeu.UtilisateurData.Navigate = back.Navigate
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	if !Header(w) {
		return
	}

	if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/search?q=%2525" + back.GetRandomLetter() + "%2525&type=playlist&limit=" + strconv.Itoa(limit) + "&offset=" + back.RandOffset()); back.Fail.Error.Status != 200 {
		fmt.Println("-----------------Erreur :-----------------", back.Fail)
		InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
		return
	}
	json.Unmarshal(back.Body, &back.Playlists)

	back.Jeu.ListPlaylist = back.Playlists
	InitTemps.Temp.ExecuteTemplate(w, "index", back.Jeu)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	if !Header(w) {
		return
	}

	Href := r.URL.Query().Get("href")
	back.Jeu.Cat = r.URL.Query().Get("type")
	if back.Body, back.Fail = back.RequestApi(Href); back.Fail.Error.Status != 200 {
		fmt.Println("-----------------Erreur :-----------------", back.Fail)
		InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
		return
	}

	switch back.Jeu.Cat {
	case "playlist":
		json.Unmarshal(back.Body, &back.PlaylistOff)
		back.PlaylistOff.Temps = back.PlaylistOff.TempsPlaylist()
		for i := 0; i < len(back.PlaylistOff.Tracks.Items); i++ {
			back.PlaylistOff.Tracks.Items[i].Track.Temps = back.Tmps(back.PlaylistOff.Tracks.Items[i].Track.DurationMs)
		}
		back.Jeu.PlaylistDetail = back.PlaylistOff
	case "artist":
		json.Unmarshal(back.Body, &back.ArtistOff)
		back.Jeu.ArtistsDetail = back.ArtistOff
	case "album":
		json.Unmarshal(back.Body, &back.AlbumOff)
		back.Jeu.AlbumsDetail = back.AlbumOff
		milsec := 0
		for i := 0; i < len(back.Jeu.AlbumsDetail.Tracks.Items); i++ {
			milsec += back.Jeu.AlbumsDetail.Tracks.Items[i].DurationMs
			back.Jeu.AlbumsDetail.Tracks.Items[i].Temps = back.Tmps(back.Jeu.AlbumsDetail.Tracks.Items[i].DurationMs)
		}
		back.Jeu.AlbumsDetail.Temps = back.Tmps(milsec)
		back.AlbumOff.Temps = back.AlbumOff.TempsAlbum()
	case "track":
		json.Unmarshal(back.Body, &back.TrackOff)
		back.Jeu.TracksDetail = back.TrackOff
		back.Jeu.TracksDetail.Temps = back.Tmps(back.Jeu.TracksDetail.DurationMs)
		back.Jeu.TracksDetail.Like = back.Like(back.Jeu.TracksDetail.Name + " " + back.Jeu.TracksDetail.Artists[0].Name)
		if back.Body, back.Fail = back.RequestApi(back.Jeu.TracksDetail.Album.Href); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
		json.Unmarshal(back.Body, &back.Jeu.AlbumsDetail)
		for i := 0; i < len(back.Jeu.AlbumsDetail.Tracks.Items); i++ {
			back.Jeu.AlbumsDetail.Tracks.Items[i].Temps = back.Tmps(back.Jeu.AlbumsDetail.Tracks.Items[i].DurationMs)
		}

	default:
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoBack(), http.StatusMovedPermanently)
		return
	}
	InitTemps.Temp.ExecuteTemplate(w, "detail", back.Jeu)
}

func Category(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}

	if !Header(w) {
		return
	}
	back.Jeu.Cat = r.URL.Query().Get("c")
	if r.URL.Query().Get("moveA") != "" {
		if back.Body, back.Fail = back.RequestApi(r.URL.Query().Get("moveA")); back.Fail.Error.Status != 200 {
			fmt.Println(r.URL.Query().Get("moveA"), "moveA")
			fmt.Println("-----------------Erreur reading requestApi 1 catgory :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
	} else if r.URL.Query().Get("moveB") != "" {
		fmt.Println(r.URL.Query().Get("moveB"), "moveB")
		if back.Body, back.Fail = back.RequestApi(r.URL.Query().Get("moveB")); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur reading requestApi 2 category :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
	} else {
		if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/search?q=%2525a%2525&type=" + back.Jeu.Cat + "&offset=0&limit=" + strconv.Itoa(limit)); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur reading requestApi 2 category :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
	}
	fmt.Println(back.Jeu.Cat)
	switch back.Jeu.Cat {
	case "playlist":
		json.Unmarshal(back.Body, &back.Playlists)
		//back.Playlists.Temps = back.TempsPlaylist(back.Playlists)
		back.Jeu.ListPlaylist = back.Playlists

	case "artist":
		json.Unmarshal(back.Body, &back.Artist)
		back.Jeu.ListArtists = back.Artist

	case "album":
		json.Unmarshal(back.Body, &back.Album)
		back.Jeu.ListAlbums = back.Album
		// for i := 0; i < len(back.Jeu.AlbumsDetail.Tracks.Items); i++ {
		// 	back.Jeu.AlbumsDetail.Tracks.Items[i].Temps = back.Tmps(back.Jeu.AlbumsDetail.Tracks.Items[i].DurationMs)
		// }
		// back.AlbumOff.Temps = back.TempsAlbum(back.AlbumOff)

	case "track":
		json.Unmarshal(back.Body, &back.Tracks)
		back.Jeu.ListTracks = back.Tracks

	default:
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		return
	}
	InitTemps.Temp.ExecuteTemplate(w, "category", back.Jeu)
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	var Max int = 100
	var Follow int
	var aime int
	if !Header(w) {
		return
	}
	for nb := 0; nb < Max; nb++ {
		if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/search?q=%2525" + r.URL.Query().Get("search") + "%2525&type=track%2Calbum%2Cartist%2Cplaylist&offset=0&limit=" + strconv.Itoa(Max)); back.Fail.Error.Status != 200 {
			if back.Fail.Error.Message == "Invalid limit" {
				Max--
				continue
			}
			fmt.Println("-----------------Erreur reading requestApi 3 Search :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
	}
	if Follow, err = strconv.Atoi(r.FormValue("f")); err != nil || Follow < 0 {
		fmt.Println("-----------------Erreur Atoi Search :-----------------", err)
		InitTemps.Temp.ExecuteTemplate(w, "search", back.Jeu)
		return
	}
	Follow=0
	r.ParseForm()
	Genre := r.Form["g"]
	json.Unmarshal(back.Body, &back.Jeu.ListAlbums)
	json.Unmarshal(back.Body, &back.Jeu.ListArtists)
	json.Unmarshal(back.Body, &back.Jeu.ListTracks)
	json.Unmarshal(back.Body, &back.Jeu.ListPlaylist)
	if Follow > 0 {
		for i, c := range back.Jeu.ListPlaylist.Playlists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi Playlist :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.PlaylistOff)
			if back.PlaylistOff.Followers.Total < Follow {
				back.Jeu.ListPlaylist.Playlists.Items = append(back.Jeu.ListPlaylist.Playlists.Items[:i], back.Jeu.ListPlaylist.Playlists.Items[i+1:]...)
			}
		}
		for i, c := range back.Jeu.ListAlbums.Albums.Items { //Filtre nul pour lui peut être voir sur une autre platforme pour récupérer les likes d'une albums
			if back.Body, back.Fail = back.RequestApi(c.Artists[0].Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi Albums :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)
			if back.ArtistOff.Followers.Total < Follow {
				back.Jeu.ListAlbums.Albums.Items = append(back.Jeu.ListAlbums.Albums.Items[:i], back.Jeu.ListAlbums.Albums.Items[i+1:]...)
			}
		}
		for i, c := range back.Jeu.ListArtists.Artists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi Artists :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)
			if back.ArtistOff.Followers.Total < Follow {
				back.Jeu.ListArtists.Artists.Items = append(back.Jeu.ListArtists.Artists.Items[:i], back.Jeu.ListArtists.Artists.Items[i+1:]...)
			}
		}
		for i, c := range back.Jeu.ListTracks.Tracks.Items {
			if aime = back.Like(c.Name + c.Artists[0].Name); aime < Follow {
				back.Jeu.ListTracks.Tracks.Items = append(back.Jeu.ListTracks.Tracks.Items[:i], back.Jeu.ListTracks.Tracks.Items[i+1:]...)
			}

		}
	}
	if len(Genre) > 0 {
		for i, c := range back.Jeu.ListPlaylist.Playlists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi 2 Playlist :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.PlaylistOff)
			for _, j := range back.PlaylistOff.Tracks.Items {
				for _, k := range j.Track.Artists {
					if !back.IsElementPresent(k.Genres, Genre) {
						back.Jeu.ListPlaylist.Playlists.Items = append(back.Jeu.ListPlaylist.Playlists.Items[:i], back.Jeu.ListPlaylist.Playlists.Items[i+1:]...)
					}
				}
			}
		}
		for i, c := range back.Jeu.ListAlbums.Albums.Items {
			if back.Body, back.Fail = back.RequestApi(c.Artists[0].Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi 2 Albums :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)
			if !back.IsElementPresent(back.ArtistOff.Genres, Genre) {
				back.Jeu.ListAlbums.Albums.Items = append(back.Jeu.ListAlbums.Albums.Items[:i], back.Jeu.ListAlbums.Albums.Items[i+1:]...)
			}
		}
		for i, c := range back.Jeu.ListArtists.Artists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi 2 Artists :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)
			if !back.IsElementPresent(back.ArtistOff.Genres, Genre) {
				back.Jeu.ListArtists.Artists.Items = append(back.Jeu.ListArtists.Artists.Items[:i], back.Jeu.ListArtists.Artists.Items[i+1:]...)
			}
		}
		for i, c := range back.Jeu.ListTracks.Tracks.Items {
			if !back.IsElementPresent(c.Artists[0].Genres, Genre) {
				back.Jeu.ListTracks.Tracks.Items = append(back.Jeu.ListTracks.Tracks.Items[:i], back.Jeu.ListTracks.Tracks.Items[i+1:]...)
			}
		}
	}
	if r.FormValue("alpha") == "ok" {
		back.Jeu.ListPlaylist.InsertionSortPlaylist()
		back.Jeu.ListAlbums.InsertionSortAlbums()
		back.Jeu.ListArtists.InsertionSortArtists()
		back.Jeu.ListTracks.InsertionSortTracks()
	}

	InitTemps.Temp.ExecuteTemplate(w, "search", back.Jeu)
}

func InitFav(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if back.LstUser, err = back.ReadJSON(); err != nil {
			fmt.Println("-----------------Error ID InitFav-----------------", err.Error())
			return
		}
		IdQuery := r.URL.Query().Get("id")
		Type := r.URL.Query().Get("type")

		if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/" + Type + "/" + IdQuery); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur reading requestApi 1 InitFav :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}

		switch Type {
		case "playlist":
			json.Unmarshal(back.Body, &back.PlaylistOff)
			back.Jeu.Utilisateur.FavorisPlaylist = append(back.Jeu.Utilisateur.FavorisPlaylist, back.PlaylistOff)
		case "artist":
			json.Unmarshal(back.Body, &back.ArtistOff)
			back.Jeu.Utilisateur.FavorisArtist = append(back.Jeu.Utilisateur.FavorisArtist, back.ArtistOff)
		case "album":
			json.Unmarshal(back.Body, &back.AlbumOff)
			back.Jeu.Utilisateur.FavorisAlbum = append(back.Jeu.Utilisateur.FavorisAlbum, back.AlbumOff)
		case "track":
			json.Unmarshal(back.Body, &back.TrackOff)
			back.Jeu.Utilisateur.FavorisTrack = append(back.Jeu.Utilisateur.FavorisTrack, back.TrackOff)
		default:
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
			return
		}

		for i, c := range back.LstUser {
			if c.Id == back.Jeu.Utilisateur.Id {
				back.LstUser[i] = back.Jeu.Utilisateur
				break
			}
		}
		back.EditJSON(back.LstUser) //Met les données dans le JSON
	}
	http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
}

func Fav(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if r.URL.Query().Get("url") == "" {
			back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
		}
		if !Header(w) {
			return
		}
		InitTemps.Temp.ExecuteTemplate(w, "fav", back.Jeu)
		return
	}
	http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
}

func Propos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	InitTemps.Temp.ExecuteTemplate(w, "propos", back.Jeu)
}

func Suppr(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if back.LstUser, err = back.ReadJSON(); err != nil {
			fmt.Println("-----------------Error ID Suppr-----------------", err.Error())
			return
		}

		queryID := r.URL.Query().Get("id")     //Récupére l'Id donné dans le Query string
		queryType := r.URL.Query().Get("type") //Récupére le type donné dans le Query string

		switch queryType {
		case "playlist":
			for i, c := range back.Jeu.Utilisateur.FavorisPlaylist {
				if c.ID == queryID {
					back.Jeu.Utilisateur.FavorisPlaylist = append(back.Jeu.Utilisateur.FavorisPlaylist[:i], back.Jeu.Utilisateur.FavorisPlaylist[i+1:]...) //Supprime de la liste des personnage
					break
				}
			}
		case "artist":
			for i, c := range back.Jeu.Utilisateur.FavorisArtist {
				if c.ID == queryID {
					back.Jeu.Utilisateur.FavorisArtist = append(back.Jeu.Utilisateur.FavorisArtist[:i], back.Jeu.Utilisateur.FavorisArtist[i+1:]...) //Supprime de la liste des personnage
					break
				}
			}
		case "album":
			for i, c := range back.Jeu.Utilisateur.FavorisAlbum {
				if c.ID == queryID {
					back.Jeu.Utilisateur.FavorisAlbum = append(back.Jeu.Utilisateur.FavorisAlbum[:i], back.Jeu.Utilisateur.FavorisAlbum[i+1:]...) //Supprime de la liste des personnage
					break
				}
			}
		case "track":
			for i, c := range back.Jeu.Utilisateur.FavorisTrack {
				if c.ID == queryID {
					back.Jeu.Utilisateur.FavorisTrack = append(back.Jeu.Utilisateur.FavorisTrack[:i], back.Jeu.Utilisateur.FavorisTrack[i+1:]...) //Supprime de la liste des personnage
					break
				}
			}
		default:
			fmt.Println("queryType false :" + queryType)
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
			return
		}

		for i, c := range back.LstUser {
			if c.Id == back.Jeu.Utilisateur.Id {
				back.LstUser[i] = back.Jeu.Utilisateur
				break
			}
		}
		back.EditJSON(back.LstUser) //Met les données dans le JSON
	}
	http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
}

func Play(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("t")
	artist := r.URL.Query().Get("a")
	hrefQuery := r.URL.Query().Get("href")
	var vid string
	if vid, err = back.Download(back.IdYtb(title + " " + artist)); err != nil {
		fmt.Println("-----------------Erreur download(-----------------")
	} else {
		if back.Body, back.Fail = back.RequestApi(hrefQuery); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur reading requestApi 1 Play :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
		json.Unmarshal(back.Body, &back.Jeu.Footer)
		back.Jeu.Footer.IDYtb = vid
		fmt.Println("You're currently playing " + title + " of " + artist)
		back.Jeu.UtilisateurData.Play = true
	}
	http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
}
