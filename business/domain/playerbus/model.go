package playerbus

import (
	"github.com/jootd/soccer-manager/business/types/age"
	"github.com/jootd/soccer-manager/business/types/position"
)

type Player struct {
	ID        int
	TeamID    int               `json:"team_id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Age       age.Age           `json:"age"`
	Country   string            `json:"country"`
	Value     int64             `json:"value"`
	Position  position.Position `json:"position"`
}

type UpdatePlayer struct {
	ID        int
	TeamID    *int     `json:"-"`
	FirstName *string  `json:"first_name"`
	LastName  *string  `json:"last_name"`
	Country   *string  `json:"country"`
	Age       *age.Age `json:"-"`
	Value     *int64   `json:"-"`
}

type CreatePlayer struct {
	TeamID    int
	FirstName string
	LastName  string
	Age       age.Age
	Country   string
	Value     int64
	Position  position.Position
}
