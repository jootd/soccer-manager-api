package playerapp

import (
	"encoding/json"

	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/types/age"
	"github.com/jootd/soccer-manager/business/types/position"
)

type UpdatePlayer struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Country   *string `json:"country"`
}

func (app UpdatePlayer) Validate() error {
	//TODO: validation
	return nil
}

func (app *UpdatePlayer) Decode(data []byte) error {
	return nil
}

func toBusUpdate(update UpdatePlayer) playerbus.UpdatePlayer {
	return playerbus.UpdatePlayer{
		FirstName: update.FirstName,
		LastName:  update.LastName,
		Country:   update.Country,
	}
}

type Player struct {
	ID        int    `json:"id"`
	TeamID    int    `json:"team_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Country   string `json:"country"`
	Value     int64  `json:"value"`
	Position  string `json:"position"`
}

func toAppPlayer(bus playerbus.Player) Player {
	return Player{
		ID:        bus.ID,
		TeamID:    bus.TeamID,
		FirstName: bus.FirstName,
		LastName:  bus.LastName,
		Age:       bus.Age.Value(),
		Country:   bus.Country,
		Value:     bus.Value,
		Position:  bus.Position.String(),
	}
}

type CreatePlayer struct {
	TeamID    int    `json:"team_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Country   string `json:"country"`
	Value     int64  `json:"value"`
	Position  string `json:"position"`
}

// Decode implements the decoder interface

func (app CreatePlayer) Decode(data []byte) error {
	return json.Unmarshal(data, &app)
}

func (app CreatePlayer) Validate() error {
	//TODO:
	return nil
}

func toBusCreatePlayer(req CreatePlayer) playerbus.CreatePlayer {
	return playerbus.CreatePlayer{
		TeamID:    req.TeamID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       age.MustParse(req.Age),
		Country:   req.Country,
		Value:     req.Value,
		Position:  position.MustParse(req.Position),
	}

}
