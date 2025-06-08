package teamapp

import "github.com/jootd/soccer-manager/business/domain/teambus"

type Update struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Budget  int64  `json:"budget"`
}

func toAppTeam(bus teambus.Team) Team {
	return Team{
		ID:      bus.ID,
		Name:    bus.Name,
		Country: bus.Country,
		Budget:  bus.Budget,
	}
}

// TODO: CreateTeam,
