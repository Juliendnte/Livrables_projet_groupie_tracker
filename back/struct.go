package back

type Site struct {
	Utilisateur     Client
	UtilisateurData ClientData

	Header Albums
	Footer TrackPrecision

	ListArtists  Artists
	ListAlbums   Albums
	ListTracks   Track
	ListPlaylist Playlist

	ArtistsDetail  ArtistPrecision
	AlbumsDetail   AlbumPrecision
	TracksDetail   TrackPrecision
	PlaylistDetail PlaylistPrecision

	MoveA string
	MoveB string
	Cat   string
}

type Artists struct { 
	Artists struct {
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Total int `json:"total"`
			} `json:"followers"`
			Genres []string `json:"genres"`
			Href   string   `json:"href"`
			ID     string   `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			Type       string `json:"type"`
		} `json:"items"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
	} `json:"artists"`
}

type Albums struct { 
	Albums struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Items    []struct {
			Href         string `json:"href"`
			AlbumType    string `json:"album_type"`
			TotalTracks  int    `json:"total_tracks"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name        string `json:"name"`
			ReleaseDate string `json:"release_date"`
			Type        string `json:"type"`
			Artists     []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"albums"`
}

type Track struct { 
	Tracks struct {
		Items []struct {
			Album struct {
				TotalTracks  int `json:"total_tracks"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name        string `json:"name"`
				ReleaseDate string `json:"release_date"`
				Artists     []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Followers struct {
					Total int `json:"total"`
				} `json:"followers"`
				Genres []string `json:"genres"`
				ID     string   `json:"id"`
				Href   string   `json:"href"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name       string `json:"name"`
				Popularity int    `json:"popularity"`
			} `json:"artists"`
			DurationMs   int `json:"duration_ms"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			ID         string `json:"id"`
			Href       string `json:"href"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			PreviewURL string `json:"preview_url"`
			Type       string `json:"type"`
		} `json:"items"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
	} `json:"tracks"`
}

type Playlist struct {
	Playlists struct {
		Items []struct {
			Description  string `json:"description"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height interface{} `json:"height"`
				URL    string      `json:"url"`
				Width  interface{} `json:"width"`
			} `json:"images"`
			Name  string `json:"name"`
			Owner struct {
				DisplayName  string `json:"display_name"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				ID string `json:"id"`
			} `json:"owner"`
			PrimaryColor interface{} `json:"primary_color"`
			Public       interface{} `json:"public"`
			Tracks       struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"tracks"`
			Type string `json:"type"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"playlists"`
}

type AlbumPrecision struct {
	AlbumType    string `json:"album_type"`
	TotalTracks  int    `json:"total_tracks"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Type        string `json:"type"`
	Artists     []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"artists"`
	Tracks struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"artists"`
			DiscNumber   int `json:"disc_number"`
			DurationMs   int `json:"duration_ms"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			PreviewURL  string `json:"preview_url"`
			Href        string `json:"href"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			Temps       string
		} `json:"items"`
	} `json:"tracks"`
	Copyrights []struct {
		Text string `json:"text"`
	} `json:"copyrights"`
	Genres     []string `json:"genres"`
	Label      string   `json:"label"`
	Popularity int      `json:"popularity"`
	Temps      string
}

type ArtistPrecision struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Genres []string `json:"genres"`
	ID     string   `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Type       string `json:"type"`
}

type PlaylistPrecision struct {
	Description  string `json:"description"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Total int `json:"total"`
		} `json:"followers"`
		ID          string `json:"id"`
		Type        string `json:"type"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Tracks struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			AddedAt string `json:"added_at"`
			Track   struct {
				Album struct {
					AlbumType    string `json:"album_type"`
					TotalTracks  int    `json:"total_tracks"`
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href   string `json:"href"`
					ID     string `json:"id"`
					Images []struct {
						URL    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"images"`
					Name        string `json:"name"`
					ReleaseDate string `json:"release_date"`
					Type        string `json:"type"`
					Artists     []struct {
						ExternalUrls struct {
							Spotify string `json:"spotify"`
						} `json:"external_urls"`
						ID   string `json:"id"`
						Name string `json:"name"`
						Type string `json:"type"`
					} `json:"artists"`
				} `json:"album"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Followers struct {
						Total int `json:"total"`
					} `json:"followers"`
					Genres []string `json:"genres"`
					ID     string   `json:"id"`
					Images []struct {
						URL    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"images"`
					Name       string `json:"name"`
					Popularity int    `json:"popularity"`
					Type       string `json:"type"`
				} `json:"artists"`
				DiscNumber   int `json:"disc_number"`
				DurationMs   int `json:"duration_ms"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				ID          string `json:"id"`
				Href        string `json:"href"`
				Name        string `json:"name"`
				Popularity  int    `json:"popularity"`
				PreviewURL  string `json:"preview_url"`
				TrackNumber int    `json:"track_number"`
				Type        string `json:"type"`
				Temps       string
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
	Type  string `json:"type"`
	Temps string
}

type TrackPrecision struct {
	Album struct {
		AlbumType    string `json:"album_type"`
		TotalTracks  int    `json:"total_tracks"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name        string `json:"name"`
		ReleaseDate string `json:"release_date"`
		Type        string `json:"type"`
		Artists     []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"artists"`
	} `json:"album"`
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Total int `json:"total"`
		} `json:"followers"`
		Genres []string `json:"genres"`
		Href   string   `json:"href"`
		ID     string   `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		Type       string `json:"type"`
	} `json:"artists"`
	DiscNumber   int `json:"disc_number"`
	DurationMs   int `json:"duration_ms"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	Preview_Url string `json:"preview_url"`
	Href        string `json:"href"`
	TrackNumber int    `json:"track_number"`
	Type        string `json:"type"`
	Like        int
	IDYtb       string
	Temps       string
}

type ClientData struct {
	Connect  bool
	Play     bool
}

type Client struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Mdp    string `json:"mdp"`
	Img    string `json:"img"`
	Letter string `json:"letter"`

	FavorisTrack    []TrackPrecision    `json:"favoristrack"`
	FavorisPlaylist []PlaylistPrecision `json:"favorisplaylist"`
	FavorisArtist   []ArtistPrecision   `json:"favorisartist"`
	FavorisAlbum    []AlbumPrecision    `json:"favorisalbum"`
}

type ErreurApi struct {
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
}

type Duree struct {
	heure int
	min   int
	sec   int
}

var Body []byte
var Fail ErreurApi

var TracksAll Track
var AlbumAll Albums
var ArtistAll Artists
var PlaylistsAll Playlist

var Tracks Track
var Album Albums
var Artist Artists
var Playlists Playlist

var TrackOff TrackPrecision
var AlbumOff AlbumPrecision
var ArtistOff ArtistPrecision
var PlaylistOff PlaylistPrecision

var LstUser []Client
var User Client
var UserData ClientData

var Jeu = Site{Utilisateur: User, UtilisateurData: UserData}
