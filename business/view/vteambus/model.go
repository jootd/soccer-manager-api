package vteambus

import (
	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/teambus"
)

type Player struct {
	ID        int    `json:"id"`
	FirstName string `json:"name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Value     int64  `json:"value"`
}

type TeamWithPlayers struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Country string   `json:"country"`
	Players []Player `json:"players"`
	Value   int64    `json:"value"`
}

func FromTeam(team teambus.Team, players []playerbus.Player) TeamWithPlayers {
	viewPlayers := make([]Player, 0, len(players))
	var teamValue int64
	for _, p := range players {
		if p.TeamID != team.ID {
			continue
		}
		viewPlayer := Player{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Age:       p.Age,
		}
		teamValue += p.Value
		viewPlayers = append(viewPlayers, viewPlayer)
	}

	return TeamWithPlayers{
		ID:      team.ID,
		Name:    team.Name,
		Country: team.Country,
		Players: viewPlayers,
		Value:   teamValue,
	}
}

func FromTeams(teams []teambus.Team, players []playerbus.Player) []TeamWithPlayers {
	result := make([]TeamWithPlayers, 0, len(teams))

	playerMap := make(map[int][]playerbus.Player)
	for _, p := range players {
		playerMap[p.TeamID] = append(playerMap[p.TeamID], p)
	}

	for _, t := range teams {
		result = append(result, FromTeam(t, playerMap[t.ID]))
	}

	return result
}
