package riot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RiotAPISettings builds the RiotAPI
type RiotAPISettings struct {
	APIKey string
	Region string
}

// Create creates a RiotAPI
func (r RiotAPISettings) Create() *RiotAPI {
	return &RiotAPI{
		RiotAPISettings: r,
		apiBase:         fmt.Sprintf("https://%s.api.pvp.net", r.Region),
		apiLol:          fmt.Sprintf("https://%s.api.pvp.net/api/lol/%s", r.Region, r.Region),
	}
}

// RiotAPI is the Riot API interface
type RiotAPI struct {
	RiotAPISettings
	apiBase string
	apiLol  string
}

type FeaturedGamesResponse struct {
	GameList []FeaturedGame `json:"gameList"`
}

type FeaturedGame struct {
	GameId       int                       `json:"gameId"`
	Participants []FeaturedGameParticipant `json:"participants"`
}

type FeaturedGameParticipant struct {
	SummonerName string `json:"summonerName"`
}

// FeaturedGames gets featured games
func (r *RiotAPI) FeaturedGames() (*FeaturedGamesResponse, error) {
	resp, err := r.fetch(
		fmt.Sprintf("%s/observer-mode/rest/featured", r.apiBase))
	if err != nil {
		return nil, err
	}
	var ret FeaturedGamesResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Game gets recent games of a summoner
func (r *RiotAPI) Game(summonerId string) (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerId))
}

type MatchResponse struct {
}

// Match gets match details
func (r *RiotAPI) Match(matchId string) (*MatchResponse, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchId))
}

type SummonerByNameResponse map[string]Summoner

type Summoner struct {
	Id int `json:"id"`
}

// SummonerByName gets multiple summoners by name
func (r *RiotAPI) SummonerByName(summonerNames []string) (SummonerByNameResponse, error) {
	resp, err := r.fetchWithKey(
		fmt.Sprintf("%s/v1.4/summoner/by-name/%s",
			r.apiLol, strings.Join(summonerNames, ",")))
	if err != nil {
		return nil, err
	}
	var ret SummonerByNameResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *RiotAPI) fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

func (r *RiotAPI) fetchWithKey(u string) (*http.Response, error) {
	return r.fetchWithKeyAndParams(u, url.Values{})
}

func (r *RiotAPI) fetchAndParams(u string, params url.Values) (*http.Response, error) {
	params["api_key"] = []string{r.APIKey}
	return r.fetch(fmt.Sprintf("%s?%s", u, params.Encode()))
}
