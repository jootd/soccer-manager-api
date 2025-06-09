package teamdb

import (
	"database/sql"

	"github.com/jootd/soccer-manager/business/domain/teambus"
)

type team struct {
	ID      int            `db:"id"`
	Name    sql.NullString `db:"name"`
	Country sql.NullString `db:"country"`
	Budget  int64          `db:"budget"`
}

// Converts from local DB model to teambus model
func toBusTeam(team team) teambus.Team {
	return teambus.Team{
		ID:      team.ID,
		Name:    team.Name.String,
		Country: team.Country.String,
		Budget:  team.Budget,
	}
}

// Converts from teambus model to local DB model
func toDBTeam(bt teambus.Team) team {
	return team{
		ID: bt.ID,
		Name: sql.NullString{
			Valid:  len(bt.Name) > 0,
			String: bt.Name,
		},
		Country: sql.NullString{
			Valid:  len(bt.Country) > 0,
			String: bt.Country,
		},
		Budget: bt.Budget,
	}
}

func toBusTeamSlice(dbTeams []team) []teambus.Team {
	teams := []teambus.Team{}
	for _, db := range dbTeams {
		teams = append(teams, toBusTeam(db))
	}
	return teams
}
