package teambus

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

const (
	InitialTeamBudget = 5_000_000 //$
)

var (
	ErrTeamNotFound = errors.New("team_not_found")
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	GetByID(ctx context.Context, id int) (Team, error)
	Query(ctx context.Context, query QueryFilter) ([]Team, error)
	Update(ctx context.Context, updates Team) error
	Create(ctx context.Context, team Team) (int, error)
}

type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	GetByID(ctx context.Context, id int) (Team, error)
	Query(ctx context.Context, query QueryFilter) ([]Team, error)
	Update(ctx context.Context, updates UpdateTeam) error
	Create(ctx context.Context, team CreateTeam) error
	AutoGenerate(ctx context.Context) (Team, error)
}

type Extension func(ExtBusiness) ExtBusiness

type Business struct {
	store Storer
	log   *zap.SugaredLogger
}

func NewTeamBus(store Storer, log *zap.SugaredLogger, extensions ...Extension) ExtBusiness {
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

func (tb *Business) GetByID(ctx context.Context, id int) (Team, error) {
	team, err := tb.store.GetByID(ctx, id)
	if err != nil {
		return Team{}, fmt.Errorf("teambus:GetByID:%w", err)
	}

	return team, nil
}

func (tb *Business) Query(ctx context.Context, query QueryFilter) ([]Team, error) {
	teams, err := tb.store.Query(ctx, query)
	if err != nil {
		return []Team{}, fmt.Errorf("teambus:Query:%w", err)
	}

	return teams, nil
}

// Generate generates random team
func (tb *Business) AutoGenerate(ctx context.Context) (Team, error) {
	newTeam := Team{
		Name:    rand.Text(),
		Country: rand.Text()[:2],
		Budget:  InitialTeamBudget,
	}
	id, err := tb.store.Create(ctx, newTeam)
	if err != nil {
		return Team{}, fmt.Errorf("teambus:AutoGenerate:%w", err)
	}

	newTeam.ID = id
	return newTeam, nil
}

func (tb *Business) Create(ctx context.Context, new CreateTeam) error {
	_, err := tb.store.Create(ctx, Team{
		Name:    new.Name,
		Country: new.Country,
		Budget:  new.Budget,
	})
	if err != nil {
		return fmt.Errorf("teambus:CreateTeam:%w", err)
	}

	return nil
}

func (tb *Business) Update(ctx context.Context, upd UpdateTeam) error {
	team := Team{
		ID: upd.ID,
	}

	if upd.Country != nil {
		team.Country = *upd.Country
	}

	if upd.Name != nil {
		team.Name = *upd.Name
	}

	err := tb.store.Update(ctx, team)
	if err != nil {
		return fmt.Errorf("teambus:Update:%w", err)
	}

	return nil
}
