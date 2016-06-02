package riot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// MatchResponse is the match response
type MatchResponse struct {
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities"`
	RawJSON               string                `json:"-"`
}

// ParticipantIdentity is the identity of a participant
type ParticipantIdentity struct {
	ParticipantId int                 `json:"participantId"`
	Player        MatchResponsePlayer `json:"player"`
}

// MatchResponsePlayer is a player of a match response
type MatchResponsePlayer struct {
	SummonerId   int    `json:"summonerId"`
	SummonerName string `json:"summonerName"`
}

// Match gets match details
func (r *API) Match(matchId string) (*MatchResponse, error) {
	resp, err := r.fetch(
		fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchId))
	var m MatchResponse
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read match response: %v", err)
	}
	err = json.Unmarshal(s, &m)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal match response: %v", err)
	}
	m.RawJSON = string(s)
	return &m, nil
}
