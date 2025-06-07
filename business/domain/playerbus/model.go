package playerbus

import (
	"github.com/jootd/soccer-manager/business/types/age"
	"github.com/jootd/soccer-manager/business/types/position"
)

type Player struct {
	ID        int
	TeamID    int
	FirstName string
	LastName  string
	Age       age.Age
	Country   string
	Value     int64
	Position  position.Position
}

type UpdatePlayer struct {
	FirstName *string
	LastName  *string
	Country   *string
	Age       *age.Age
	Value     *int64
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
