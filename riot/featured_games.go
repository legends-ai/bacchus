package riot

import (
	"encoding/json"
	"fmt"
)

// FeaturedGamesResponse is the featured games response
type FeaturedGamesResponse struct {
	GameList []FeaturedGame `json:"gameList"`
}

// FeaturedGame is a featured game
type FeaturedGame struct {
	GameID       int                       `json:"gameId"`
	Participants []FeaturedGameParticipant `json:"participants"`
}

// FeaturedGameParticipant is a participant of a featured game
type FeaturedGameParticipant struct {
	SummonerName string `json:"summonerName"`
}

// FeaturedGames gets featured games
func (r *API) FeaturedGames() (*FeaturedGamesResponse, error) {
	resp, err := r.fetch(
		fmt.Sprintf("%s/observer-mode/rest/featured", r.apiBase))
	if err != nil {
		return nil, err
	}
	var ret FeaturedGamesResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
