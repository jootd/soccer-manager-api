package playerdb

import "github.com/jootd/soccer-manager/business/domain/playerbus"

type player struct {
	Id        int
	FirstName string  `db:"first_name"`
	LastName  string  `db:"last_name"`
	Country   string  `db:"country"`
	Value     float64 `db:"value"`
	TeamId    int     `db:"team_id"`
}

func toPlayerBus(player player) playerbus.Player {
	return playerbus.Player{
		ID:        player.Id,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		Country:   player.Country,
		Value:     player.Value,
		TeamId:    player.TeamId,
	}
}

func toDBPlayer(bus playerbus.Player) player {
	return player{
		Id:        bus.ID,
		FirstName: bus.FirstName,
		LastName:  bus.LastName,
		Country:   bus.Country,
		Value:     bus.Value,
		TeamId:    bus.TeamId,
	}
}
