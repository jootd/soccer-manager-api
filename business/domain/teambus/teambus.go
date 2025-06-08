package teambus

import (
	"context"
	"fmt"

	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

const (
	InitialTeamBudget = 5_000_000 //$
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Query(ctx context.Context, query QueryFilter) ([]Team, error)
	Update(ctx context.Context, updates Team) error
	Create(ctx context.Context, team Team) error
}

type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Query(ctx context.Context, query QueryFilter) ([]Team, error)
	Update(ctx context.Context, team Team, updates UpdateTeam) error
	Create(ctx context.Context, team CreateTeam) error
}

type Extension func(ExtBusiness) ExtBusiness

type Business struct {
	store Storer
	log   *zap.Logger
}

func NewTeamBus(store Storer, log *zap.Logger, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		store: store,
		log:   log,
	})

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

func (tb *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := tb.store.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := &Business{
		log:   tb.log,
		store: storer,
	}

	return bus, nil

}

func (tb *Business) Query(ctx context.Context, query QueryFilter) ([]Team, error) {
	teams, err := tb.store.Query(ctx, query)
	if err != nil {
		return []Team{}, fmt.Errorf("teambus:Query:%w", err)
	}

	return teams, nil
}

// Generate generates random team
func (tb *Business) Generate(ctx context.Context) {

}

func (tb *Business) Create(ctx context.Context, new CreateTeam) error {
	err := tb.store.Create(ctx, Team{
		Name:    new.Name,
		Country: new.Country,
		Budget:  new.Budget,
	})
	if err != nil {
		return fmt.Errorf("teambus:CreateTeam:%w", err)
	}

	return nil

}

func (tb *Business) Update(ctx context.Context, team Team, updates UpdateTeam) error {
	//TODO : update logic
	err := tb.store.Update(ctx, Team{})
	if err != nil {
		return fmt.Errorf("teambus:Update:%w", err)
	}

	return nil
}
