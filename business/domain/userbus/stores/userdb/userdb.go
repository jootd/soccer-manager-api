package userdb

import (
	"context"
	"log"

	"github.com/jootd/soccer-manager/business/domain/userbus"
)

type Store struct {
	log *log.Logger
	//sqlx
}

func NewStore(logger *log.Logger) *Store {
	return &Store{
		log: logger,
	}
}

func (us *Store) GetUser(ctx context.Context, username string) (userbus.User, error) {
	user := user{}

	return toBusUser(user), nil
}

func (us *Store) Create(ctx context.Context, username string, passHash string) (userbus.User, error) {
	newUser := user{}

	return toBusUser(newUser), nil
}

func (us *Store) Update(ctx context.Context, username string, teamId int) (userbus.User, error) {
	user := userbus.User{}

	return user, nil
}
