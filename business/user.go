package business

import (
	"context"
	"fmt"
)

type User struct {
	Username string
	Password string
	TeamId   int
}

type UserStorer interface {
	GetUser(ctx context.Context, username string) (User, error)
	CreateUser(ctx context.Context, username string, passHash string) (User, error)
	UpdateUser(ctx context.Context, username string, teamId int) (User, error)
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
	user, err := ub.store.GetUser(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("bus:GetUser:%w", err)
	}

	return user, nil

}

func (ub *UserBus) CreateUser(ctx context.Context, username string, passHash string) (User, error) {
	user, err := ub.store.CreateUser(ctx, username, passHash)
	if err != nil {
		return User{}, fmt.Errorf("bus:CreateUser:%w", err)
	}

	return user, nil
}
func (ub *UserBus) UpdateUser(ctx context.Context, username string, teamId int) (User, error) {
	user, err := ub.store.UpdateUser(ctx, username, teamId)
	if err != nil {
		return User{}, fmt.Errorf("bus:UpdateUser:%w", err)
	}

	return user, nil
}
