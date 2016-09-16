package db

import (
	"github.com/asunaio/bacchus/models"
	"github.com/gocql/gocql"
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
func (m *MatchesDAO) Exists(id models.MatchID) (bool, error) {
	var count int
	if err := m.Session.Query(hasMatchQuery, id.String()).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}

// Insert inserts a match to Cassandra.
func (a *MatchesDAO) Insert(m *models.Match) error {
	return a.Session.Query(
		insertMatchQuery, m.ID.String(), m.ID.Region,
		m.Body, m.Rank.ToNumber(), m.Patch,
	).Exec()
}
