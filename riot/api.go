package riot

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// APISettings builds the API
type APISettings struct {
	Region string
}

// Create creates a API
func (r APISettings) Create() *API {
	base := os.Getenv("GRAGAS_RIOT_BASE")
	if base == "" {
		base = "http://localhost:3006"
	}
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("riot-region", r.Region)
	return client.Do(req)
}

func (r *API) fetchWithParams(u string, params url.Values) (*http.Response, error) {
	return r.fetch(fmt.Sprintf("%s?%s", u, params.Encode()))
}
