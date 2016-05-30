package riot

import (
	"fmt"
	"io/ioutil"
)

type MatchResponse struct {
	RawJSON string `json:"-"`
}

// Match gets match details
func (r *API) Match(matchId string) (*MatchResponse, error) {
	resp, err := r.fetchWithKey(
		fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchId))
	fmt.Println(fmt.Sprintf("%s/v2.2/match/%s", r.apiLol, matchId))
	defer resp.Body.Close()
	var m MatchResponse
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read match response: %v", err)
	}
	m.RawJSON = string(s)
	return &m, nil
}
