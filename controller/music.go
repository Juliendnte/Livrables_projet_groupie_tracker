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
var search string

func Header(w http.ResponseWriter, detail bool) bool {
	if back.Jeu.UtilisateurData.Fav {
		if detail {
			back.Jeu.UtilisateurData.FavorisAlbums = []back.AlbumPrecision{}
			back.Jeu.UtilisateurData.FavorisArtists = []back.ArtistPrecision{}
			back.Jeu.UtilisateurData.FavorisPlaylist = []back.PlaylistPrecision{}
			back.Jeu.UtilisateurData.FavorisTracks = []back.TrackPrecision{}

			if back.LstUser, err = back.ReadJSON(); err != nil {
				fmt.Println("-----------------Error ID InitFav-----------------", err.Error())
				InitTemps.Temp.ExecuteTemplate(w, "404", back.Jeu)
				return false
			}
			for _, c := range back.LstUser {
				if c.Id == back.Jeu.Utilisateur.Id {
					for _, v := range c.FavorisAlbum {
						if back.Body, back.Fail = back.RequestApi(v); back.Fail.Error.Status != 200 {
							fmt.Println("-----------------Erreur :-----------------", back.Fail)
							InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
							return false
						}
						album := back.AlbumPrecision{}
						json.Unmarshal(back.Body, &album)
						back.Jeu.UtilisateurData.FavorisAlbums = append(back.Jeu.UtilisateurData.FavorisAlbums, album)
					}
					for _, v := range c.FavorisArtist {
						if back.Body, back.Fail = back.RequestApi(v); back.Fail.Error.Status != 200 {
							fmt.Println("-----------------Erreur :-----------------", back.Fail)
							InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
							return false
						}
						Artist := back.ArtistPrecision{}
						json.Unmarshal(back.Body, &Artist)
						back.Jeu.UtilisateurData.FavorisArtists = append(back.Jeu.UtilisateurData.FavorisArtists, Artist)
					}
					for _, v := range c.FavorisPlaylist {
						if back.Body, back.Fail = back.RequestApi(v); back.Fail.Error.Status != 200 {
							fmt.Println("-----------------Erreur :-----------------", back.Fail)
							InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
							return false
						}
						Playlist := back.PlaylistPrecision{}
						json.Unmarshal(back.Body, &Playlist)
						back.Jeu.UtilisateurData.FavorisPlaylist = append(back.Jeu.UtilisateurData.FavorisPlaylist, Playlist)
					}
					for _, v := range c.FavorisTrack {
						if back.Body, back.Fail = back.RequestApi(v); back.Fail.Error.Status != 200 {
							fmt.Println("-----------------Erreur :-----------------", back.Fail)
							InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
							return false
						}
						track := back.TrackPrecision{}
						json.Unmarshal(back.Body, &track)
						back.Jeu.UtilisateurData.FavorisTracks = append(back.Jeu.UtilisateurData.FavorisTracks, track)
					}
				}
			}
			back.Jeu.Header.Albums.Href = ""
		}

	} else if back.Jeu.Header.Albums.Href == "" {
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
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	if !Header(w, false) {
		return
	}

	if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/search?q=%2525" + back.GetRandomLetter() + "%2525&type=playlist&limit=" + strconv.Itoa(limit) + "&offset=" + back.RandOffset()); back.Fail.Error.Status != 200 {
		fmt.Println("-----------------Erreur :-----------------", back.Fail)
		InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
		return
	}
	json.Unmarshal(back.Body, &back.Jeu.ListPlaylist)

	InitTemps.Temp.ExecuteTemplate(w, "index", back.Jeu)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	if !Header(w, false) {
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
			back.PlaylistOff.Tracks.Items[i].AddedAt = back.PlaylistOff.Tracks.Items[i].AddedAt[:10]
		}
		back.Jeu.PlaylistDetail = back.PlaylistOff
	case "artist":
		json.Unmarshal(back.Body, &back.Jeu.ArtistsDetail)
		back.Jeu.ArtistsDetail.Followers.Totalstr = back.FormatAbo(back.Jeu.ArtistsDetail.Followers.Total)
		if back.Body, back.Fail = back.RequestApi(Href + "/albums"); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
		json.Unmarshal(back.Body, &back.Jeu.AlbumListAPI)
		if back.Body, back.Fail = back.RequestApi(Href + "/top-tracks?limit=5"); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
		json.Unmarshal(back.Body, &back.Jeu.TrackListAPI)
		if len(back.Jeu.TrackListAPI.Tracks) > 4 {
			back.Jeu.TrackListAPI.Tracks = back.Jeu.TrackListAPI.Tracks[:5]
		}
		for i := 0; i < len(back.Jeu.TrackListAPI.Tracks); i++ {
			back.Jeu.TrackListAPI.Tracks[i].Temps = back.Tmps(back.Jeu.TrackListAPI.Tracks[i].DurationMs)
			back.Jeu.TrackListAPI.Tracks[i].Like = back.Like(back.Jeu.TrackListAPI.Tracks[i].Name + " " + back.Jeu.TrackListAPI.Tracks[i].Artists[0].Name)
		}
	case "album":
		json.Unmarshal(back.Body, &back.Jeu.AlbumsDetail)
		milsec := 0
		for i := 0; i < len(back.Jeu.AlbumsDetail.Tracks.Items); i++ {
			milsec += back.Jeu.AlbumsDetail.Tracks.Items[i].DurationMs
			back.Jeu.AlbumsDetail.Tracks.Items[i].Temps = back.Tmps(back.Jeu.AlbumsDetail.Tracks.Items[i].DurationMs)
		}
		back.Jeu.AlbumsDetail.Temps = back.Tmps(milsec)
		back.Jeu.AlbumsDetail.Temps = back.Jeu.AlbumsDetail.TempsAlbum()
	case "track":
		json.Unmarshal(back.Body, &back.Jeu.TracksDetail)
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

	if !Header(w, false) {
		return
	}
	back.Jeu.Cat = r.URL.Query().Get("c")
	if r.URL.Query().Get("moveA") != "" {
		if back.Body, back.Fail = back.RequestApi(r.URL.Query().Get("moveA")); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur reading requestApi 1 catgory :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
	} else if r.URL.Query().Get("moveB") != "" {
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

	switch back.Jeu.Cat {
	case "playlist":
		json.Unmarshal(back.Body, &back.Jeu.ListPlaylist)
	case "artist":
		json.Unmarshal(back.Body, &back.Jeu.ListArtists)
		for i, c := range back.Jeu.ListArtists.Artists.Items {
			if len(c.Images) == 0 {
				back.Jeu.ListArtists.Artists.Items[i].Images = append(back.Jeu.ListArtists.Artists.Items[i].Images, back.Images{URL: "/static/img/nopdp.png"})
			}
		}
	case "album":
		json.Unmarshal(back.Body, &back.Jeu.ListAlbums)
	case "track":
		json.Unmarshal(back.Body, &back.Jeu.ListTracks)

	default:
		if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
			http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		}
		return
	}
	InitTemps.Temp.ExecuteTemplate(w, "category", back.Jeu)
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	var Follow int
	var aime int
	back.Jeu.ListAlbums = back.Albums{}
	back.Jeu.ListArtists = back.Artists{}
	back.Jeu.ListPlaylist = back.Playlist{}
	back.Jeu.ListTracks = back.Track{}

	if !Header(w, false) {
		return
	}

	if len(back.Jeu.AllGenres.Genres) < 1 {
		if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/recommendations/available-genre-seeds"); back.Fail.Error.Status != 200 {
			fmt.Println("-----------------Erreur reading requestApi 1 search :-----------------", back.Fail)
			InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
			return
		}
		json.Unmarshal(back.Body, &back.Jeu.AllGenres)
	}
	var strFollow string

	if r.URL.Query().Get("search") != "" {
		search = r.URL.Query().Get("search")
	} else if search == "" {
		InitTemps.Temp.ExecuteTemplate(w, "search", back.Jeu)
		return
	}

	if back.Body, back.Fail = back.RequestApi("https://api.spotify.com/v1/search?q=" + search + "&type=track%2Calbum%2Cartist%2Cplaylist&offset=0&limit=10"); back.Fail.Error.Status != 200 {
		fmt.Println("-----------------Erreur reading requestApi 2 Search :-----------------", back.Fail)
		InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
		return
	}

	json.Unmarshal(back.Body, &back.Album)
	json.Unmarshal(back.Body, &back.Artist)
	json.Unmarshal(back.Body, &back.Playlists)
	json.Unmarshal(back.Body, &back.Tracks)

	if strFollow = r.FormValue("follower"); strFollow != "" {
		if Follow, err = strconv.Atoi(strFollow); err != nil || Follow < 0 {
			fmt.Println("-----------------Erreur Atoi Search :-----------------", err)
			if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
				http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, "/index", http.StatusMovedPermanently)
			}
		}
	} else {
		Follow = 0
	}

	r.ParseForm()
	Genre := r.Form["genre"]
	if Follow > 0 {
		for i, c := range back.Playlists.Playlists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi Playlist :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.PlaylistOff)
			if back.PlaylistOff.Followers.Total >= Follow {
				back.Jeu.ListPlaylist.Playlists.Items = append(back.Jeu.ListPlaylist.Playlists.Items, back.Playlists.Playlists.Items[i])
			}
		}
		for i, c := range back.Album.Albums.Items { //Filtre nul pour lui peut être voir sur une autre platforme pour récupérer les likes d'un albums (deezer)
			if back.Body, back.Fail = back.RequestApi(c.Artists[0].Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi Albums :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)
			fmt.Println(back.ArtistOff.Followers.Total)
			if back.ArtistOff.Followers.Total >= Follow {
				back.Jeu.ListAlbums.Albums.Items = append(back.Jeu.ListAlbums.Albums.Items, back.Album.Albums.Items[i])
			}
		}
		for i, c := range back.Artist.Artists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi Artists :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)

			if back.ArtistOff.Followers.Total >= Follow {
				fmt.Println(c.Name)
				back.Jeu.ListArtists.Artists.Items = append(back.Jeu.ListArtists.Artists.Items, back.Artist.Artists.Items[i])
			}
		}
		for i, c := range back.Tracks.Tracks.Items {
			if len(c.Artists) > 0 {
				if aime, _ = strconv.Atoi(back.Like(c.Name + " " + c.Artists[0].Name)); aime >= Follow {
					back.Jeu.ListTracks.Tracks.Items = append(back.Jeu.ListTracks.Tracks.Items, back.Tracks.Tracks.Items[i])
				} else {
					fmt.Println("son ", aime)
				}
			}
		}
	} else {
		back.Jeu.ListPlaylist = back.Playlists
		back.Jeu.ListAlbums = back.Album
		back.Jeu.ListTracks = back.Tracks
		back.Jeu.ListArtists = back.Artist
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
					if back.IsElementPresent(k.Genres, Genre) {
						back.Jeu.ListPlaylist.Playlists.Items = append(back.Jeu.ListPlaylist.Playlists.Items, back.Playlists.Playlists.Items[i])
					}
				}
			}
		}
		for i, c := range back.Jeu.ListAlbums.Albums.Items {
			if len(c.Artists) > 0 {
				if back.Body, back.Fail = back.RequestApi(c.Artists[0].Href); back.Fail.Error.Status != 200 {
					fmt.Println("-----------------Erreur reading requestApi 2 Albums :-----------------", back.Fail)
					InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
					return
				}
				json.Unmarshal(back.Body, &back.ArtistOff)
				if back.IsElementPresent(back.ArtistOff.Genres, Genre) {
					back.Jeu.ListAlbums.Albums.Items = append(back.Jeu.ListAlbums.Albums.Items, back.Album.Albums.Items[i])
				}
			}
		}
		for i, c := range back.Artist.Artists.Items {
			if back.Body, back.Fail = back.RequestApi(c.Href); back.Fail.Error.Status != 200 {
				fmt.Println("-----------------Erreur reading requestApi 2 Artists :-----------------", back.Fail)
				InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
				return
			}
			json.Unmarshal(back.Body, &back.ArtistOff)
			if back.IsElementPresent(back.ArtistOff.Genres, Genre) {
				back.Jeu.ListArtists.Artists.Items = append(back.Jeu.ListArtists.Artists.Items, back.Artist.Artists.Items[i])
			}
		}
		for i, c := range back.Tracks.Tracks.Items {
			if len(c.Artists) > 0 {
				if back.IsElementPresent(c.Artists[0].Genres, Genre) {
					back.Jeu.ListTracks.Tracks.Items = append(back.Jeu.ListTracks.Tracks.Items, back.Tracks.Tracks.Items[i])
				}
			}
		}
	}
	if Follow == 0 && len(Genre) == 0 {
		json.Unmarshal(back.Body, &back.Jeu.ListAlbums)
		json.Unmarshal(back.Body, &back.Jeu.ListArtists)
		json.Unmarshal(back.Body, &back.Jeu.ListPlaylist)
		json.Unmarshal(back.Body, &back.Jeu.ListTracks)
	}
	fmt.Println(r.FormValue("alpha"))
	if r.FormValue("alpha") == "ok" {
		back.Jeu.ListPlaylist.InsertionSortPlaylist()
		back.Jeu.ListAlbums.InsertionSortAlbums()
		back.Jeu.ListArtists.InsertionSortArtists()
		back.Jeu.ListTracks.InsertionSortTracks()
	}
	for i := 0; i < len(back.Jeu.ListTracks.Tracks.Items); i++ {
		back.Jeu.ListTracks.Tracks.Items[i].Temps = back.Tmps(back.Jeu.ListTracks.Tracks.Items[i].DurationMs)
	}
	for i, c := range back.Jeu.ListArtists.Artists.Items {
		if len(c.Images) == 0 {
			back.Jeu.ListArtists.Artists.Items[i].Images = append(back.Jeu.ListArtists.Artists.Items[i].Images, back.Images{URL: "/static/img/nopdp.png"})
		}
	}
	back.Jeu.ListTracks.Tracks.Total =len(back.Jeu.ListTracks.Tracks.Items)
	back.Jeu.ListPlaylist.Playlists.Total =len(back.Jeu.ListPlaylist.Playlists.Items)
	back.Jeu.ListAlbums.Albums.Total =len(back.Jeu.ListAlbums.Albums.Items)
	back.Jeu.ListArtists.Artists.Total =len(back.Jeu.ListArtists.Artists.Items)

	fmt.Println(back.Jeu.ListPlaylist, back.Jeu.ListAlbums, back.Jeu.ListArtists, back.Jeu.ListTracks)
	InitTemps.Temp.ExecuteTemplate(w, "search", back.Jeu)
}

func Propos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("url") == "" {
		back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
	}
	InitTemps.Temp.ExecuteTemplate(w, "propos", back.Jeu)
}

func Fav(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if r.URL.Query().Get("url") == "" {
			back.Jeu.UtilisateurData.Navigate.VisitPage(r.URL.String())
		}

		if back.LstUser, err = back.ReadJSON(); err != nil {
			fmt.Println("-----------------Error ID InitFav-----------------", err.Error())
			InitTemps.Temp.ExecuteTemplate(w, "404", back.Jeu)
			return
		}
		clear(back.Jeu.ListAlbumsDetail)
		clear(back.Jeu.ListArtistsDetail)
		clear(back.Jeu.ListTracksDetail)
		clear(back.Jeu.ListPlaylistDetail)

		for i, c := range back.LstUser {
			if c.Id == back.Jeu.Utilisateur.Id {
				back.LstUser[i] = back.Jeu.Utilisateur
				for _, s := range c.FavorisAlbum {
					if back.Body, back.Fail = back.RequestApi(s); back.Fail.Error.Status != 200 {
						fmt.Println("-----------------Erreur reading requestApi 1 InitFav :-----------------", back.Fail)
						InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
						return
					}
					json.Unmarshal(back.Body, &back.AlbumOff)
					back.Jeu.ListAlbumsDetail = append(back.Jeu.ListAlbumsDetail, back.AlbumOff)
				}

				for _, s := range c.FavorisArtist {
					if back.Body, back.Fail = back.RequestApi(s); back.Fail.Error.Status != 200 {
						fmt.Println("-----------------Erreur reading requestApi 1 InitFav :-----------------", back.Fail)
						InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
						return
					}
					json.Unmarshal(back.Body, &back.ArtistOff)
					back.Jeu.ListArtistsDetail = append(back.Jeu.ListArtistsDetail, back.ArtistOff)
				}

				for _, s := range c.FavorisTrack {
					if back.Body, back.Fail = back.RequestApi(s); back.Fail.Error.Status != 200 {
						fmt.Println("-----------------Erreur reading requestApi 1 InitFav :-----------------", back.Fail)
						InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
						return
					}
					json.Unmarshal(back.Body, &back.TrackOff)
					back.Jeu.ListTracksDetail = append(back.Jeu.ListTracksDetail, back.TrackOff)
					for i := 0; i < len(back.Jeu.ListTracksDetail); i++ {
						back.Jeu.ListTracksDetail[i].Temps = back.Tmps(back.Jeu.ListTracksDetail[i].DurationMs)
					}
				}

				for _, s := range c.FavorisPlaylist {
					if back.Body, back.Fail = back.RequestApi(s); back.Fail.Error.Status != 200 {
						fmt.Println("-----------------Erreur reading requestApi 1 InitFav :-----------------", back.Fail)
						InitTemps.Temp.ExecuteTemplate(w, strconv.Itoa(back.Fail.Error.Status), back.Jeu)
						return
					}
					json.Unmarshal(back.Body, &back.PlaylistOff)
					back.Jeu.ListPlaylistDetail = append(back.Jeu.ListPlaylistDetail, back.PlaylistOff)
				}
				break
			}
		}
		InitTemps.Temp.ExecuteTemplate(w, "fav", back.Jeu)
		return
	}
	if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
}

func InitFav(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if back.LstUser, err = back.ReadJSON(); err != nil {
			fmt.Println("-----------------Error ID InitFav-----------------", err.Error())
			InitTemps.Temp.ExecuteTemplate(w, "404", back.Jeu)
			return
		}
		Href := r.URL.Query().Get("href")
		Type := r.URL.Query().Get("type")

		switch Type {
		case "playlist":
			if !back.IsInList(back.Jeu.Utilisateur.FavorisPlaylist, Href) {
				back.Jeu.UtilisateurData.Fav = true
				back.Jeu.Utilisateur.FavorisPlaylist = append(back.Jeu.Utilisateur.FavorisPlaylist, Href)
			}
		case "artist":
			if !back.IsInList(back.Jeu.Utilisateur.FavorisArtist, Href) {
				back.Jeu.UtilisateurData.Fav = true
				back.Jeu.Utilisateur.FavorisArtist = append(back.Jeu.Utilisateur.FavorisArtist, Href)
			}
		case "album":
			if !back.IsInList(back.Jeu.Utilisateur.FavorisAlbum, Href) {
				back.Jeu.UtilisateurData.Fav = true
				back.Jeu.Utilisateur.FavorisAlbum = append(back.Jeu.Utilisateur.FavorisAlbum, Href)
			}
		case "track":
			if !back.IsInList(back.Jeu.Utilisateur.FavorisTrack, Href) {
				back.Jeu.Utilisateur.FavorisTrack = append(back.Jeu.Utilisateur.FavorisTrack, Href)
				back.Jeu.UtilisateurData.Fav = true
			}
		default:
			if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
				http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, "/index", http.StatusMovedPermanently)
			}
			return
		}

		for i, c := range back.LstUser {
			if c.Id == back.Jeu.Utilisateur.Id {
				back.LstUser[i] = back.Jeu.Utilisateur
				break
			}
		}
		back.EditJSON(back.LstUser) //Met les données dans le JSON
		if !Header(w, true) {
			return
		}
	}
	if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
}

func Suppr(w http.ResponseWriter, r *http.Request) {
	if back.Jeu.UtilisateurData.Connect {
		if back.LstUser, err = back.ReadJSON(); err != nil {
			fmt.Println("-----------------Error ID Suppr-----------------", err.Error())
			return
		}

		queryHref := r.URL.Query().Get("href") //Récupére l'Id donné dans le Query string
		queryType := r.URL.Query().Get("type") //Récupére le type donné dans le Query string
		switch queryType {
		case "playlist":
			for i, c := range back.Jeu.Utilisateur.FavorisPlaylist {
				if c == queryHref {
					back.Jeu.Utilisateur.FavorisPlaylist = append(back.Jeu.Utilisateur.FavorisPlaylist[:i], back.Jeu.Utilisateur.FavorisPlaylist[i+1:]...)
					back.Jeu.UtilisateurData.FavorisPlaylist = append(back.Jeu.UtilisateurData.FavorisPlaylist[:i], back.Jeu.UtilisateurData.FavorisPlaylist[i+1:]...)
					break
				}
			}
		case "artist":
			for i, c := range back.Jeu.Utilisateur.FavorisArtist {
				if c == queryHref {
					back.Jeu.Utilisateur.FavorisArtist = append(back.Jeu.Utilisateur.FavorisArtist[:i], back.Jeu.Utilisateur.FavorisArtist[i+1:]...)
					back.Jeu.UtilisateurData.FavorisArtists = append(back.Jeu.UtilisateurData.FavorisArtists[:i], back.Jeu.UtilisateurData.FavorisArtists[i+1:]...)
					break
				}
			}
		case "album":
			for i, c := range back.Jeu.Utilisateur.FavorisAlbum {
				if c == queryHref {
					back.Jeu.Utilisateur.FavorisAlbum = append(back.Jeu.Utilisateur.FavorisAlbum[:i], back.Jeu.Utilisateur.FavorisAlbum[i+1:]...)
					back.Jeu.UtilisateurData.FavorisAlbums = append(back.Jeu.UtilisateurData.FavorisAlbums[:i], back.Jeu.UtilisateurData.FavorisAlbums[i+1:]...)

					break
				}
			}
		case "track":
			for i, c := range back.Jeu.Utilisateur.FavorisTrack {
				if c == queryHref {
					back.Jeu.Utilisateur.FavorisTrack = append(back.Jeu.Utilisateur.FavorisTrack[:i], back.Jeu.Utilisateur.FavorisTrack[i+1:]...)
					back.Jeu.UtilisateurData.FavorisTracks = append(back.Jeu.UtilisateurData.FavorisTracks[:i], back.Jeu.UtilisateurData.FavorisTracks[i+1:]...)
					break
				}
			}
		default:
			fmt.Println("queryType false :" + queryType)
			if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
				http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, "/index", http.StatusMovedPermanently)
			}
			return
		}

		for i, c := range back.LstUser {
			if c.Id == back.Jeu.Utilisateur.Id {
				back.LstUser[i] = back.Jeu.Utilisateur
				break
			}
		}

		if len(back.Jeu.Utilisateur.FavorisAlbum) == 0 && len(back.Jeu.Utilisateur.FavorisArtist) == 0 && len(back.Jeu.Utilisateur.FavorisPlaylist) == 0 && len(back.Jeu.Utilisateur.FavorisTrack) == 0 {
			back.Jeu.UtilisateurData.Fav = false
		}

		if !Header(w, true) {
			return
		}
		back.EditJSON(back.LstUser) //Met les données dans le JSON
	}
	if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
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
	if len(back.Jeu.UtilisateurData.Navigate.History) > 0 {
		http.Redirect(w, r, back.Jeu.UtilisateurData.Navigate.GoAcutal(), http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
}
