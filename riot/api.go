package riot

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	envRiot = "RIOT_BASE"
)

// Create creates a API
func Create(region string) *API {
	base := os.Getenv(envRiot)
	if base == "" {
		base = "http://localhost:3006"
	}
	return &API{
		Region:  region,
		apiBase: base,
		apiLol:  fmt.Sprintf("%s/api/lol/%s", base, region),
	}
}

// API is the Riot API interface
type API struct {
	Region  string
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
