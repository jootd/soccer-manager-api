package playerdb

import (
	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/types/age"
	"github.com/jootd/soccer-manager/business/types/position"
)

type player struct {
	ID        int    `db:"id"`
	TeamID    int    `db:"team_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int    `db:"age"`
	Country   string `db:"country"`
	Value     int64  `db:"value"`
	Position  string `db:"position"`
}

func toPlayerBus(player player) playerbus.Player {
	return playerbus.Player{
		ID:        player.ID,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		Age:       age.MustParse(player.Age),
		Country:   player.Country,
		Value:     player.Value,
		TeamID:    player.TeamID,
		Position:  position.MustParse(player.Position),
	}
}

func toPlayerBusSlice(players []player) []playerbus.Player {
	var bus []playerbus.Player

	for _, p := range players {
		bus = append(bus, toPlayerBus(p))
	}

	return bus
}

func toDBPlayer(bus playerbus.Player) player {
	return player{
		ID:        bus.ID,
		FirstName: bus.FirstName,
		LastName:  bus.LastName,
		Age:       bus.Age.Value(),
		Country:   bus.Country,
		Value:     bus.Value,
		TeamID:    bus.TeamID,
		Position:  bus.Position.String(),
	}
}

func toDBPlayerSlice(bus []playerbus.Player) []player {
	var db []player
	for _, p := range bus {
		db = append(db, toDBPlayer(p))
	}
	return db
}
