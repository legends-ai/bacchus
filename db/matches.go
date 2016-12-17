package db

import (
	"github.com/gocql/gocql"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
)

const (
	hasMatchQuery    = `SELECT COUNT(*) FROM match_set WHERE id = ?`
	insertMatchQuery = `INSERT INTO match_set (id) VALUES (?)`
)

// MatchesDAO is a matches DAO.
type MatchesDAO struct {
	// Session is the session to the Athena cluster.
	Session *gocql.Session `inject:"t"`
}

// Exists returns true if a match exists.
func (m *MatchesDAO) Exists(id *apb.MatchId) (bool, error) {
	var count int
	if err := m.Session.Query(hasMatchQuery, models.StringifyMatchId(id)).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}

// Insert inserts a match id to Cassandra.
func (m *MatchesDAO) Insert(id *apb.MatchId) error {
	return m.Session.Query(
		insertMatchQuery, models.StringifyMatchId(id),
	).Exec()
}
