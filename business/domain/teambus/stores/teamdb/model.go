package teamdb

import "github.com/jootd/soccer-manager/business/domain/teambus"

type team struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	Country string `db:"country"`
	Budget  int64  `db:"budget"`
}

// Converts from local DB model to teambus model
func toBusTeam(team team) teambus.Team {
	return teambus.Team{
		ID:      team.ID,
		Name:    team.Name,
		Country: team.Country,
		Budget:  team.Budget,
	}
}

// Converts from teambus model to local DB model
func toDBTeam(bt teambus.Team) team {
	return team{
		ID:      bt.ID,
		Name:    bt.Name,
		Country: bt.Country,
		Budget:  bt.Budget,
	}
}

func toBusTeamSlice(dbTeams []team) []teambus.Team {
	teams := []teambus.Team{}
	for _, db := range dbTeams {
		teams = append(teams, toBusTeam(db))
	}
	return teams
}
