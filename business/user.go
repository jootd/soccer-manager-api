package business

import (
	"context"
	"errors"
)

type User struct {
	Username string
	Password string
	TeamId   int
}

type UserStorer interface {
	GetUser(ctx context.Context, username string) (User, bool)
	CreateUser(ctx context.Context, username string, passHash string) (User, bool)
	UpdateUser(ctx context.Context, username string, teamId int) (User, bool)
}

type UserBus struct {
	store UserStorer
}

func NewUserBus(store UserStorer) *UserBus {
	return &UserBus{
		store: store,
	}
}

func (ub *UserBus) GetUser(ctx context.Context, username string) (User, error) {
	user, ok := ub.store.GetUser(ctx, username)
	if !ok {
		return User{}, errors.New("")
	}

	return user, nil

}

func (ub *UserBus) CreateUser(ctx context.Context, username string, passHash string) (User, error) {
	user, ok := ub.store.CreateUser(ctx, username, passHash)
	if !ok {
		return User{}, errors.New("")
	}

	return user, nil
}
func (ub *UserBus) UpdateUser(ctx context.Context, username string, teamId int) (User, error) {
	user, ok := ub.store.UpdateUser(ctx, username, teamId)
	if !ok {
		return User{}, errors.New("")
	}

	return user, nil
}
