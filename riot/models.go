package riot

type Summoner struct {
	AccountID     string `json:"accountId"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	ID            string `json:"id"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	SummonerLevel int    `json:"summonerLevel"`
}

type LeagueEntry struct {
	LeagueID     string `json:"leagueId"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	HotStreak    bool   `json:"hotStreak"`
	Veteran      bool   `json:"veteran"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
}

type Match struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

type MatchMetadata struct {
	DataVersion  string   `json:"data_version"`
	MatchID      string   `json:"match_id"`
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	GameDatetime    int64         `json:"game_datetime"`
	GameLength      float64       `json:"game_length"`
	GameVersion     string        `json:"game_version"`
	Participants    []Participant `json:"participants"`
	QueueID         int           `json:"queue_id"`
	TftGameType     string        `json:"tft_game_type"`
	TftSetCoreNeme  string        `json:"tft_set_core_name"`
	TftSetNumber    int           `json:"tft_set_number"`
}

type Participant struct {
	Augments              []string    `json:"augments"`
	Companion             Companion   `json:"companion"`
	GoldLeft              int         `json:"gold_left"`
	LastRound             int         `json:"last_round"`
	Level                 int         `json:"level"`
	Placement             int         `json:"placement"`
	PlayersEliminated     int         `json:"players_eliminated"`
	PUUID                 string      `json:"puuid"`
	TimeEliminated        float64     `json:"time_eliminated"`
	TotalDamageToPlayers  int         `json:"total_damage_to_players"`
	Traits                []Trait     `json:"traits"`
	Units                 []Unit      `json:"units"`
}

type Companion struct {
	ContentID   string `json:"content_ID"`
	ItemID      int    `json:"item_ID"`
	SkinID      int    `json:"skin_ID"`
	Species     string `json:"species"`
}

type Trait struct {
	Name          string `json:"name"`
	NumUnits      int    `json:"num_units"`
	Style         int    `json:"style"`
	TierCurrent   int    `json:"tier_current"`
	TierTotal     int    `json:"tier_total"`
}

type Unit struct {
	Items       []int    `json:"items"`
	CharacterID string   `json:"character_id"`
	Chosen      string   `json:"chosen"`
	Name        string   `json:"name"`
	Rarity      int      `json:"rarity"`
	Tier        int      `json:"tier"`
}

type RiotAPIError struct {
	Status Status `json:"status"`
}

type Status struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}