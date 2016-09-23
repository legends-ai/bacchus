package db

import (
	"github.com/gocql/gocql"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
)

const (
	hasMatchQuery    = `SELECT COUNT(*) FROM matches WHERE id = ?`
	insertMatchQuery = `INSERT INTO matches (id, region, body, rank, patch) VALUES (?, ?, ?, ?, ?)`
)

// MatchesDAO is a matches DAO.
type MatchesDAO struct {
	// Session is the session to the Athena cluster.
	Session *gocql.Session `inject:"t"`
}

// Exists returns true if a match exists.
func (m *MatchesDAO) Exists(id *apb.MatchId) (bool, error) {
	var count int
	if err := m.Session.Query(hasMatchQuery, id.String()).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}

// Insert inserts a match to Cassandra.
func (a *MatchesDAO) Insert(m *apb.RawMatch) error {
	return a.Session.Query(
		insertMatchQuery, models.StringifyMatchId(m.Id),
		m.Id.Region, m.Body, models.RankToNumber(m.Rank),
		m.Patch,
	).Exec()
}
