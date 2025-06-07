package playerdb

import "github.com/jootd/soccer-manager/business/domain/playerbus"

type player struct {
	Id        int
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int    `db:"age"`
	Country   string `db:"country"`
	Value     int64  `db:"value"`
	TeamID    int    `db:"team_id"`
}

func toPlayerBus(player player) playerbus.Player {
	return playerbus.Player{
		ID:        player.Id,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		Age:       player.Age,
		Country:   player.Country,
		Value:     player.Value,
		TeamID:    player.TeamID,
	}
}

func toDBPlayer(bus playerbus.Player) player {
	return player{
		Id:        bus.ID,
		FirstName: bus.FirstName,
		LastName:  bus.LastName,
		Age:       bus.Age,
		Country:   bus.Country,
		Value:     bus.Value,
		TeamID:    bus.TeamID,
	}
}
