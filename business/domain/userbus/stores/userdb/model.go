package userdb

import (
	"time"

	"github.com/jootd/soccer-manager/business/domain/userbus"
)

type user struct {
	Username    string    `db:"username"`
	Password    string    `db:"password_hash"`
	TeamID      int       `db:"team_id"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toBusUser(db user) userbus.User {
	return userbus.User{
		Username:    db.Username,
		Password:    db.Password,
		TeamID:      db.TeamID,
		DateCreated: db.DateCreated,
		DateUpdated: db.DateUpdated,
	}
}

func toDBUser(bus userbus.User) user {
	return user{
		Username:    bus.Username,
		Password:    bus.Password,
		TeamID:      bus.TeamID,
		DateCreated: bus.DateCreated,
		DateUpdated: bus.DateUpdated,
	}
}
