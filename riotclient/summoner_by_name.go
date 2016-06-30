package riotclient

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SummonerByNameResponse is the summoner by name response
type SummonerByNameResponse map[string]*Summoner

// Summoner is a summoner
type Summoner struct {
	ID int `json:"id"`
}

// SummonerByName gets multiple summoners by name
func (r *API) SummonerByName(summonerNames []string) (SummonerByNameResponse, error) {
	resp, err := r.fetch(
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
