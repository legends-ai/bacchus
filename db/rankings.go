package db

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/models"
)

const (
	rankingsQuery      = `SELECT rankings FROM rankings WHERE id = ?`
	insertRankingQuery = `INSERT INTO rankings (id, rankings) VALUES (?, ?)`
	updateRankingQuery = `UPDATE rankings SET rankings = rankings + ? WHERE id = ?`
)

// RankingsDAO is a rankings DAO.
type RankingsDAO struct {
	// Session is the session to the Athena cluster.
	Session *gocql.Session `inject:"t"`
}

// Get grabs all rankings of a summoner.
func (a *RankingsDAO) Get(id models.SummonerID) (*models.RankingList, error) {
	var rankings []struct {
		Time time.Time
		Rank int64
	}
	if err := a.Session.Query(rankingsQuery, id.String()).Scan(&rankings); err != nil && err != gocql.ErrNotFound {
		return nil, err
	}
	var ret []*models.Ranking
	for _, ranking := range rankings {
		ret = append(ret, &models.Ranking{
			Time: ranking.Time,
			Rank: models.RankFromNumber(uint32(ranking.Rank)),
		})
	}
	return models.NewRankingList(ret), nil
}

// Insert stores an Athena ranking row for a new summoner.
func (a *RankingsDAO) Insert(id models.SummonerID, r models.Ranking) error {
	return a.Session.Query(insertRankingQuery, id.String(), r.UDTSet()).Exec()
}

// Update updates the Athena ranking of the given summoner with the given ranking.
func (a *RankingsDAO) Update(id models.SummonerID, r models.Ranking) error {
	return a.Session.Query(updateRankingQuery, r.UDTSet(), id.String()).Exec()
}