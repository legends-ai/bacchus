package db

import (
	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/models"
)

const (
	rankingsQuery      = `SELECT rankings FROM rankings WHERE id = ?`
	insertRankingQuery = `INSERT INTO rankings (id, rankings, rank) VALUES (?, ?, ?)`
	updateRankingQuery = `UPDATE rankings SET rankings = rankings + ? AND rank = ? WHERE id = ?`
	aboveRankQuery     = `SELECT id FROM rankings WHERE rank >= ? LIMIT ? ALLOW FILTERING`
)

// RankingsDAO is a rankings DAO.
type RankingsDAO struct {
	// Session is the session to the Athena cluster.
	Session *gocql.Session `inject:"t"`
}

// Get grabs all rankings of a summoner.
func (a *RankingsDAO) Get(id models.SummonerID) (*models.RankingList, error) {
	var rankings []models.RankingUDT
	if err := a.Session.Query(rankingsQuery, id.String()).Scan(&rankings); err != nil && err != gocql.ErrNotFound {
		return nil, err
	}
	var ret []*models.Ranking
	for _, ranking := range rankings {
		ret = append(ret, ranking.ToRanking())
	}
	return models.NewRankingList(ret), nil
}

// AboveRank gets all summoner ids above a given rank with a limit.
func (r *RankingsDAO) AboveRank(rank models.Rank, limit int) ([]models.SummonerID, error) {
	var ret []models.SummonerID
	if err := r.Session.Query(aboveRankQuery, rank.ToNumber(), limit).Scan(&ret); err != nil && err != gocql.ErrNotFound {
		return nil, err
	}
	return ret, nil
}

// Insert stores an Athena ranking row for a new summoner.
func (a *RankingsDAO) Insert(id models.SummonerID, r models.Ranking) error {
	return a.Session.Query(insertRankingQuery, id.String(), r.UDTSet(), r.Rank.ToNumber()).Exec()
}

// Update updates the Athena ranking of the given summoner with the given ranking.
func (a *RankingsDAO) Update(id models.SummonerID, r models.Ranking) error {
	return a.Session.Query(updateRankingQuery, r.UDTSet(), r.Rank.ToNumber(), id.String()).Exec()
}

// Upsert updates or inserts a ranking based on if it exists or not.
func (d *RankingsDAO) Upsert(id models.SummonerID, r models.Ranking, exists bool) error {
	if exists {
		return d.Update(id, r)
	} else {
		return d.Insert(id, r)
	}
}
