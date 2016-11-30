package db

import (
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
)

const (
	hasMatchQuery    = `SELECT COUNT(*) FROM matches_serialized WHERE id = ?`
	insertMatchQuery = `INSERT INTO matches_serialized (id, region, rank, patch, data) VALUES (?, ?, ?, ?, ?)`
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

// Insert inserts a match to Cassandra.
func (a *MatchesDAO) Insert(m *apb.BacchusData_RawMatch) error {
	data, err := proto.Marshal(m.Data)
	if err != nil {
		return err
	}

	return a.Session.Query(
		insertMatchQuery, models.StringifyMatchId(m.Id),
		m.Id.Region.String(), models.RankToNumber(m.Rank), m.Patch, data,
	).Exec()
}
