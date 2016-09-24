package models

import (
	"fmt"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

func StringifySummonerId(id *apb.SummonerId) string {
	return fmt.Sprintf("%s/%d", id.Region, id.Id)
}

func StringifyMatchId(id *apb.MatchId) string {
	return fmt.Sprintf("%s/%d", id.Region, id.Id)
}
