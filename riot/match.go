package riot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

// MatchResponse is the match response
type MatchResponse struct {
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities"`
	MatchCreation         int64                 `json:"matchCreation"`
	MatchVersion          string                `json:"matchVersion"`
	QueueType             string                `json:"queueType"`
	RawJSON               string                `json:"-"`
}

// Time returns the time of this match
func (r *MatchResponse) Time() time.Time {
	return time.Unix(r.MatchCreation/1000, r.MatchCreation%1000*1E6)
}

// ParticipantIdentity is the identity of a participant
type ParticipantIdentity struct {
	ParticipantID int                 `json:"participantId"`
	Player        MatchResponsePlayer `json:"player"`
}

// MatchResponsePlayer is a player of a match response
type MatchResponsePlayer struct {
	SummonerID   int    `json:"summonerId"`
	SummonerName string `json:"summonerName"`
}

// Match gets match details
func (r *API) Match(matchID string) (*MatchResponse, error) {
	resp, err := r.fetchWithParams(
		fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchID), url.Values{"includeTimeline": []string{"true"}})
	if err != nil {
		return nil, err
	}
	var m MatchResponse
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read match response: %v", err)
	}

	if err = json.Unmarshal(s, &m); err != nil {
		return nil, fmt.Errorf("Could not unmarshal match response: %v", err)
	}

	m.RawJSON = string(s)
	return &m, nil
}
