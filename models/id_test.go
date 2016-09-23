package models

import "testing"

func TestStringifySummonerId(t *testing.T) {
	for _, test := range []struct {
		Id   *apb.SummonerId
		Want string
	}{
		{
			Id: &apb.SummonerId{
				Region: "na",
				ID:     1738,
			},
			Want: "na/1738",
		},
	} {
		got := StringifySummonerId(test.Id)
		if got != test.Want {
			t.Errorf("Got %v, want %v", got, test.Want)
		}
	}
}

func TestStringifyMatchId(t *testing.T) {
	for _, test := range []struct {
		Id   *apb.MatchId
		Want string
	}{
		{
			Id: &apb.MatchId{
				Region: "na",
				ID:     1738,
			},
			Want: "na/1738",
		},
	} {
		got := StringifyMatchId(test.Id)
		if got != test.Want {
			t.Errorf("Got %v, want %v", got, test.Want)
		}
	}
}
