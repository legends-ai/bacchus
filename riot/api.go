package riot

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/simplyianm/keypool"
)

const (
	envRiot      = "RIOT_BASE"
	apiKeyParam  = "api_key"
	riotBaseTpl  = "https://%s.api.pvp.net"
	regionHeader = "riot-region"
)

// Client stores API clients.
type Client struct {
	Keys      *keypool.Keypool `inject:"t"`
	clients   map[string]*API
	clientsMu sync.RWMutex
}

// New creates a new Client.
func New() *Client {
	return &Client{
		clients: map[string]*API{},
	}
}

// Region gets an API client for the given region.
func (rc *Client) Region(region string) *API {
	rc.clientsMu.RLock()
	inst, ok := rc.clients[region]
	rc.clientsMu.RUnlock()
	if !ok {
		base := fmt.Sprintf(riotBaseTpl, region)
		inst = &API{
			Region:  region,
			apiBase: base,
			apiLol:  fmt.Sprintf("%s/api/lol/%s", base, region),
			rc:      rc,
		}
		rc.clientsMu.Lock()
		rc.clients[region] = inst
		rc.clientsMu.Unlock()
	}
	return inst
}

// API is the Riot API interface
type API struct {
	Region  string
	apiBase string
	apiLol  string
	rc      *Client
}

// fetchWithParams fetches a path with the given parameters.
func (r *API) fetchWithParams(path string, params url.Values) (*http.Response, error) {
	key := r.rc.Keys.Fetch().Return()
	params.Set(apiKeyParam, key)
	url := fmt.Sprintf("%s?%s", path, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// fetch fetches a path via GET request.
func (r *API) fetch(path string) (*http.Response, error) {
	return r.fetchWithParams(path, url.Values{})
}
