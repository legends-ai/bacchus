package riotclient

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/simplyianm/keypool"
)

const (
	envRiot      = "RIOT_BASE"
	apiKeyParam  = "api_key"
	riotBaseTpl  = "https://%s.api.pvp.net"
	regionHeader = "riot-region"
)

// RiotClient stores clients.
type RiotClient struct {
	Keys    *keypool.Keypool `inject:"t"`
	clients map[string]*API
}

// New creates a new RiotClient.
func New() *RiotClient {
	return &RiotClient{
		clients: map[string]*API{},
	}
}

// Region gets an API client for the given region.
func (rc *RiotClient) Region(region string) *API {
	inst, ok := rc.clients[region]
	if !ok {
		base := fmt.Sprintf(riotBaseTpl, region)
		inst = &API{
			Region:  region,
			apiBase: base,
			apiLol:  fmt.Sprintf("%s/api/lol/%s", base, region),
			rc:      rc,
		}
		rc.clients[region] = inst
	}
	return inst
}

// API is the Riot API interface
type API struct {
	Region  string
	apiBase string
	apiLol  string
	rc      *RiotClient
}

// fetchWithParams fetches a url with the given parameters.
func (r *API) fetchWithParams(u string, params url.Values) (*http.Response, error) {
	key := r.rc.Keys.Fetch().Return()
	params.Set(apiKeyParam, key)
	return r.fetch(fmt.Sprintf("%s?%s", u, params.Encode()))
}

// fetch fetches a URL via GET request. Do not use.
func (r *API) fetch(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
