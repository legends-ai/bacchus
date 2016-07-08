package models

// Match represents a match.
type Match struct {
	ID    MatchID
	Body  string
	Patch string
	Rank  Rank
}
