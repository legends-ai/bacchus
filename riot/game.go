package riot

import (
	"fmt"
	"io/ioutil"
)

type GameResponse struct {
	RawJSON string
}

// Game gets recent games of a summoner
func (r *API) Game(summonerId string) (*GameResponse, error) {
	resp, err := r.fetchWithKey(
		fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerId))
	fmt.Println(fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerId))
	defer resp.Body.Close()
	var g GameResponse
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read game response: %v", err)
	}
	g.RawJSON = string(s)
	return &g, nil
}
