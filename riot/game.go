package riot

import (
	"encoding/json"
	"fmt"
)

type GameResponse struct {
	SummonerId int    `json:"summonerId"`
	Games      []Game `json:"games"`
}

type Game struct {
	GameId int `json:"gameId"`
}

// Game gets recent games of a summoner
func (r *API) Game(summonerId string) (*GameResponse, error) {
	resp, err := r.fetch(
		fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerId))
	var g GameResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
