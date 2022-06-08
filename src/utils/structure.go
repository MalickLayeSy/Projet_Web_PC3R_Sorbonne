package utils

import "time"

type MatchJSON []struct {
	Winner struct {
		Acronym    string    `json:"acronym"`
		ID         int       `json:"id"`
		ImageURL   string    `json:"image_url"`
		Location   string    `json:"location"`
		ModifiedAt time.Time `json:"modified_at"`
		Name       string    `json:"name"`
		Slug       string    `json:"slug"`
	} `json:"winner"`
	Opponents []struct { //Taille 2
		Opponent struct {
			Acronym    string    `json:"acronym"`
			ID         int       `json:"id"`
			ImageURL   string    `json:"image_url"`
			Location   string    `json:"location"`
			ModifiedAt time.Time `json:"modified_at"`
			Name       string    `json:"name"`
			Slug       string    `json:"slug"`
		} `json:"opponent"`
		Type string `json:"type"`
	} `json:"opponents"`
	WinnerID      int    `json:"winner_id"`
	Name          string `json:"name"`
	NumberOfGames int    `json:"number_of_games"`
	Live          struct {
		OpensAt   time.Time `json:"opens_at"`
		Supported bool      `json:"supported"`
		URL       string    `json:"url"`
	} `json:"live"`
	DetailedStats     bool      `json:"detailed_stats"`
	Forfeit           bool      `json:"forfeit"`
	EndAt             time.Time `json:"end_at"`
	Draw              bool      `json:"draw"`
	OfficialStreamURL string    `json:"official_stream_url"`
	Tournament        struct {
		BeginAt       time.Time   `json:"begin_at"`
		EndAt         interface{} `json:"end_at"`
		ID            int         `json:"id"`
		LeagueID      int         `json:"league_id"`
		LiveSupported bool        `json:"live_supported"`
		ModifiedAt    time.Time   `json:"modified_at"`
		Name          string      `json:"name"`
		Prizepool     interface{} `json:"prizepool"`
		SerieID       int         `json:"serie_id"`
		Slug          string      `json:"slug"`
		WinnerID      interface{} `json:"winner_id"`
		WinnerType    string      `json:"winner_type"`
	} `json:"tournament"`
	BeginAt             time.Time `json:"begin_at"`
	OriginalScheduledAt time.Time `json:"original_scheduled_at"`
	Streams             struct {
		English struct {
			EmbedURL string `json:"embed_url"`
			RawURL   string `json:"raw_url"`
		} `json:"english"`
		Official struct {
			EmbedURL string `json:"embed_url"`
			RawURL   string `json:"raw_url"`
		} `json:"official"`
		Russian struct {
			EmbedURL interface{} `json:"embed_url"`
			RawURL   string      `json:"raw_url"`
		} `json:"russian"`
	} `json:"streams"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Serie       struct {
		BeginAt     time.Time   `json:"begin_at"`
		Description interface{} `json:"description"`
		EndAt       interface{} `json:"end_at"`
		FullName    string      `json:"full_name"`
		ID          int         `json:"id"`
		LeagueID    int         `json:"league_id"`
		ModifiedAt  time.Time   `json:"modified_at"`
		Name        string      `json:"name"`
		Season      string      `json:"season"`
		Slug        string      `json:"slug"`
		Tier        string      `json:"tier"`
		WinnerID    interface{} `json:"winner_id"`
		WinnerType  interface{} `json:"winner_type"`
		Year        int         `json:"year"`
	} `json:"serie"`
	TournamentID int `json:"tournament_id"`
	Results      []struct {
		Score  int `json:"score"`
		TeamID int `json:"team_id"`
	} `json:"results"`
	Status    string `json:"status"`
	MatchType string `json:"match_type"`
	League    struct {
		ID         int         `json:"id"`
		ImageURL   string      `json:"image_url"`
		ModifiedAt time.Time   `json:"modified_at"`
		Name       string      `json:"name"`
		Slug       string      `json:"slug"`
		URL        interface{} `json:"url"`
	} `json:"league"`
	Videogame struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"videogame"`
	GameAdvantage    interface{} `json:"game_advantage"`
	Slug             string      `json:"slug"`
	ID               int         `json:"id"`
	VideogameVersion struct {
		Current bool   `json:"current"`
		Name    string `json:"name"`
	} `json:"videogame_version"`
	ModifiedAt   time.Time `json:"modified_at"`
	Rescheduled  bool      `json:"rescheduled"`
	LeagueID     int       `json:"league_id"`
	LiveEmbedURL string    `json:"live_embed_url"`
	Games        []struct {
		BeginAt       time.Time   `json:"begin_at"`
		Complete      bool        `json:"complete"`
		DetailedStats bool        `json:"detailed_stats"`
		EndAt         time.Time   `json:"end_at"`
		Finished      bool        `json:"finished"`
		Forfeit       bool        `json:"forfeit"`
		ID            int         `json:"id"`
		Length        int         `json:"length"`
		MatchID       int         `json:"match_id"`
		Position      int         `json:"position"`
		Status        string      `json:"status"`
		VideoURL      interface{} `json:"video_url"`
		Winner        struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		} `json:"winner"`
		WinnerType string `json:"winner_type"`
	} `json:"games"`
	SerieID int `json:"serie_id"`
}
