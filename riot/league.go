package riot

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	QueueSolo5x5                   = "RANKED_SOLO_5x5"
	QueuePremade5x5                = "RANKED_PREMADE_5x5"
	QueueTeam3x3                   = "RANKED_TEAM_3x3"
	QueueTeam5x5                   = "RANKED_TEAM_5x5"
	QueueTeamBuilderDraftRanked5x5 = "TEAM_BUILDER_DRAFT_RANKED_5x5"
)

// LeagueResponse is the league response
type LeagueResponse map[string][]*LeagueDto

// LeagueDto contains league information.
type LeagueDto struct {
	Name    string            `json:"name"`
	Queue   string            `json:"queue"`
	Tier    string            `json:"tier"`
	Entries []*LeagueEntryDto `json:"entries"`
}

// LeagueEntryDto contains league participant information representing a summoner or team.
type LeagueEntryDto struct {
	PlayerOrTeamID   string `json:"playerOrTeamId"`
	PlayerOrTeamName string `json:"playerOrTeamName"`
	Division         string `json:"division"`
}

// League gets a league
func (r *API) League(summonerIds []string) (LeagueResponse, error) {
	idsStr := strings.Join(summonerIds, ",")
	resp, err := r.fetch(
		fmt.Sprintf("%s/v2.5/league/by-summoner/%s", r.apiLol, idsStr))
	if err != nil {
		return nil, err
	}
	ret := LeagueResponse{}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {

		return nil, err
	}
	return ret, nil
}
