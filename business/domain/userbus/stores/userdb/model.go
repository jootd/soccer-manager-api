package userdb

import "github.com/jootd/soccer-manager/business/domain/userbus"

type user struct {
	Username string
	Password string
	TeamId   int
}

func toBusUser(db user) userbus.User {
	return userbus.User{
		Username: db.Username,
		Password: db.Password,
		TeamID:   db.TeamId,
	}
}

func toDBUser(bus userbus.User) user {
	return user{
		Username: bus.Username,
		Password: bus.Password,
		TeamId:   bus.TeamID,
	}
}
