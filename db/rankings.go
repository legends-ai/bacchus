package db

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/models"
)

const (
	rankingsQuery      = `SELECT time, rank FROM rankings WHERE id = ?`
	insertRankingQuery = `INSERT INTO rankings (id, time, rank) VALUES (?, ?, ?)`
	aboveRankQuery     = `SELECT id FROM rankings WHERE rank >= ? LIMIT ? ALLOW FILTERING`
)

// RankingsDAO is a rankings DAO.
type RankingsDAO struct {
	// Session is the session to the Athena cluster.
	Session *gocql.Session `inject:"t"`
}

// Get grabs all rankings of a summoner.
func (a *RankingsDAO) Get(id models.SummonerID) (*models.RankingList, error) {
	var rankings []*models.Ranking
	iter := a.Session.Query(rankingsQuery, id.String()).Iter()
	var t time.Time
	var r models.Rank
	for iter.Scan(&t, &r) {
		rankings = append(rankings, &models.Ranking{
			ID:   id,
			Time: t,
			Rank: r,
		})
	}
	return models.NewRankingList(rankings), nil
}

// AboveRank gets all summoner ids above a given rank with a limit.
func (r *RankingsDAO) AboveRank(rank models.Rank, limit int) ([]models.SummonerID, error) {
	var ret []models.SummonerID
	it := r.Session.Query(aboveRankQuery, rank.ToNumber(), limit).Iter()
	var cur string
	for it.Scan(&cur) {
		id, err := models.SummonerIDFromString(cur)
		if err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}
	return ret, nil
}

// Insert stores an Athena ranking row for a summoner.
func (a *RankingsDAO) Insert(r models.Ranking) error {
	return a.Session.Query(insertRankingQuery, r.ID.String(), r.Time, r.Rank.ToNumber()).Exec()
}
