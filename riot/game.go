package riot

import (
	"encoding/json"
	"fmt"
)

// GameResponse is the game response
type GameResponse struct {
	SummonerID int    `json:"summonerId"`
	Games      []Game `json:"games"`
}

// Game is a game
type Game struct {
	GameID int `json:"gameId"`
}

// Game gets recent games of a summoner
func (r *API) Game(summonerID string) (*GameResponse, error) {
	resp, err := r.fetch(
		fmt.Sprintf("%s/v1.3/game/by-summoner/%s/recent", r.apiLol, summonerID))
	var g GameResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
