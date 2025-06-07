package userbus

import (
	"context"
	"fmt"
)

type Storer interface {
	Get(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, username string, passHash string) (User, error)
	Update(ctx context.Context, username string, teamId int) (User, error)
}

type Business struct {
	store Storer
}

func NewUserBus(store Storer) *Business {
	return &Business{
		store: store,
	}
}

func (ub *Business) Get(ctx context.Context, username string) (User, error) {
	user, err := ub.store.Get(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("bus:Get:%w", err)
	}

	return user, nil

}

func (ub *Business) Create(ctx context.Context, username string, passHash string) (User, error) {
	user, err := ub.store.Create(ctx, username, passHash)
	if err != nil {
		return User{}, fmt.Errorf("bus:Create:%w", err)
	}

	return user, nil
}
func (ub *Business) Update(ctx context.Context, username string, teamId int) (User, error) {
	user, err := ub.store.Update(ctx, username, teamId)
	if err != nil {
		return User{}, fmt.Errorf("bus:Update:%w", err)
	}

	return user, nil
}
