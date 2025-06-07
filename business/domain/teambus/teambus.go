package teambus

import (
	"context"
	"crypto/rand"
	"fmt"
)

const (
	InitialTeamBudget = 5_000_000 //$
)

type Storer interface {
	Query(ctx context.Context, query QueryFilter) ([]Team, error)
	Update(ctx context.Context, updates UpdateTeam) (Team, error)
	Create(ctx context.Context, team CreateTeam) (Team, error)
}

type Business struct {
	store Storer
}

func NewTeamBus(store Storer) *Business {
	return &Business{
		store: store,
	}
}

func (tb *Business) Query(ctx context.Context, query QueryFilter) ([]Team, error) {
	teams, err := tb.store.Query(ctx, query)
	if err != nil {
		return []Team{}, fmt.Errorf("teambus:Query:%w", err)
	}

	return teams, nil
}

func (tb *Business) Create(ctx context.Context) (Team, error) {
	team, err := tb.store.Create(ctx, CreateTeam{
		Name:    "team" + rand.Text(),
		Country: "country" + rand.Text(),
		Budget:  InitialTeamBudget,
	})

	if err != nil {
		return Team{}, fmt.Errorf("teambus:CreateTeam:%w", err)
	}
	return team, nil

}

func (tb *Business) Update(ctx context.Context, updates UpdateTeam) (Team, error) {
	team, err := tb.store.Update(ctx, updates)
	if err != nil {
		return Team{}, fmt.Errorf("teambus:Update:%w", err)
	}

	return team, nil
}
