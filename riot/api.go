package riot

import (
	"fmt"
	"net/http"
	"net/url"
)

// APISettings builds the API
type APISettings struct {
	APIKey string
	Region string
}

// Create creates a API
func (r APISettings) Create() *API {
	return &API{
		APISettings: r,
		apiBase:     fmt.Sprintf("https://%s.api.pvp.net", r.Region),
		apiLol:      fmt.Sprintf("https://%s.api.pvp.net/api/lol/%s", r.Region, r.Region),
	}
}

// API is the Riot API interface
type API struct {
	APISettings
	apiBase string
	apiLol  string
}

// Game gets recent games of a summoner
func (r *API) Game(summonerId string) (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerId))
}

// Match gets match details
func (r *API) Match(matchId string) (*http.Response, error) {
	return r.fetchWithKey(
		fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchId))
}

func (r *API) fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

func (r *API) fetchWithKey(u string) (*http.Response, error) {
	return r.fetchWithKeyAndParams(u, url.Values{})
}

func (r *API) fetchWithKeyAndParams(u string, params url.Values) (*http.Response, error) {
	params["api_key"] = []string{r.APIKey}
	return r.fetch(fmt.Sprintf("%s?%s", u, params.Encode()))
}
