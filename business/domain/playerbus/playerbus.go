package playerbus

import (
	"context"
	"fmt"
)

type Storer interface {
	Query(ctx context.Context, query QueryFilter) ([]Player, error)
	Update(ctx context.Context, player UpdatePlayer) (Player, error)
}

type Business struct {
	store Storer
}

func NewPlayerBus(store Storer) *Business {
	return &Business{
		store: store,
	}
}

func (up *Business) Update(ctx context.Context, player UpdatePlayer) (Player, error) {
	updated, err := up.store.Update(ctx, player)
	if err != nil {
		return Player{}, fmt.Errorf("playerbus:Update:%w", err)

	}

	return updated, nil
}

func (up *Business) Query(ctx context.Context, filter QueryFilter) ([]Player, error) {
	players, err := up.store.Query(ctx, filter)
	if err != nil {
		return []Player{}, fmt.Errorf("playerbus:Query:%w", err)
	}
	return players, nil
}
