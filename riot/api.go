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
	base := fmt.Sprintf("https://%s.api.pvp.net", r.Region)
	return &API{
		APISettings: r,
		apiBase:     base,
		apiLol:      fmt.Sprintf("%s/api/lol/%s", base, r.Region),
	}
}

// API is the Riot API interface
type API struct {
	APISettings
	apiBase string
	apiLol  string
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
