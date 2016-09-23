package db

import (
	"fmt"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"
)

const (
	rankingsQuery      = `SELECT ranking FROM rankings WHERE id = ?`
	insertRankingQuery = `INSERT INTO rankings (id, rank, ranking) VALUES (?, ?, ?)`
	aboveRankQuery     = `SELECT ranking FROM rankings WHERE rank >= ? LIMIT ? ALLOW FILTERING`
)

// RankingsDAO is a rankings DAO.
type RankingsDAO struct {
	// Session is the session to the Athena cluster.
	Session *gocql.Session `inject:"t"`
}

// Get grabs all rankings of a summoner.
func (a *RankingsDAO) Get(id *apb.SummonerId) (*apb.Ranking, error) {
	key := models.StringifySummonerId(id)

	var rawRanking []byte
	if err := a.Session.Query(rankingsQuery, key).Scan(&rawRanking); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching ranking from Cassandra: %v", err)
	}

	var ranking apb.Ranking
	if err := proto.Unmarshal(rawRanking, &ranking); err != nil {
		return nil, fmt.Errorf("error unmarshaling ranking: %v", err)
	}

	return &ranking, nil
}

// AboveRank gets all summoner ids above a given rank with a limit.
func (r *RankingsDAO) AboveRank(rank *apb.Rank, limit int) ([]*apb.SummonerId, error) {
	var ret []*apb.SummonerId
	it := r.Session.Query(aboveRankQuery, models.RankToNumber(rank), limit).Iter()

	var cur []byte
	for it.Scan(&cur) {
		var ranking apb.Ranking
		if err := proto.Unmarshal(cur, &ranking); err != nil {
			return nil, fmt.Errorf("error unmarshaling ranking: %v", err)
		}
		ret = append(ret, ranking.Summoner)
	}
	return ret, nil
}

// Insert stores an Athena ranking row for a summoner.
func (a *RankingsDAO) Insert(r *apb.Ranking) error {
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	return a.Session.Query(
		insertRankingQuery,
		models.StringifySummonerId(r.Summoner),
		models.RankToNumber(r.Rank), data,
	).Exec()
}
