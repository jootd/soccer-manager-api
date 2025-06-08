package teamadapter

import (
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
)

func toTransferTeamInfo(t teambus.Team) transferbus.TeamInfo {
	return transferbus.TeamInfo{
		ID:     t.ID,
		Budget: t.Budget,
	}
}
