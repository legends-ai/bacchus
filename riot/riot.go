package riot

import (
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
		apiBase:         fmt.Sprintf("http://%s.api.pvp.net", r.Region),
		apiLol:          fmt.Sprintf("http://%s.api.pvp.net/api/lol/%s", r.Region, r.Region),
	}
}

// RiotAPI is the Riot API interface
type RiotAPI struct {
	RiotAPISettings
	apiBase string
	apiLol  string
}

// FeaturedGames gets featured games
func (r *RiotAPI) FeaturedGames() (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/observer-mode/rest/featured", r.apiBase))
}

// Game gets recent games of a summoner
func (r *RiotAPI) Game(summonerId string) (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerId))
}

// Match gets match details
func (r *RiotAPI) Match(matchId string) (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchId))
}

// SummonerByName gets multiple summoners by name
func (r *RiotAPI) SummonerByName(summonerNames []string) (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v1.4/summoner/by-name/%s",
			r.apiLol, strings.Join(summonerNames, ",")))
}

func (r *RiotAPI) fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

func (r *RiotAPI) fetchWithKey(u string) (*http.Response, error) {
	return r.fetchWithKeyAndParams(u, url.Values{})
}

func (r *RiotAPI) fetchWithKeyAndParams(u string, params url.Values) (*http.Response, error) {
	params["api_key"] = []string{r.APIKey}
	return r.fetch(fmt.Sprintf("%s?%s", u, params.Encode()))
}
