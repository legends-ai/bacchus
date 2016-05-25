package riot

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SummonerByNameResponse map[string]*Summoner

type Summoner struct {
	Id int `json:"id"`
}

// SummonerByName gets multiple summoners by name
func (r *API) SummonerByName(summonerNames []string) (SummonerByNameResponse, error) {
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
