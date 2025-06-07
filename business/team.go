package business

import (
	"context"
	"crypto/rand"
	"fmt"
)

type Team struct {
	ID      int
	Name    string
	Country string
}

type UpdateTeam struct {
	Id      int
	Name    *string
	Country *string
}

type CreateTeam struct {
	Name    string
	Country string
}

type QueryTeam struct {
	ID      *int
	Name    *string
	Country *string
}

type TeamStorer interface {
	GetTeamsBy(ctx context.Context, query QueryTeam) ([]Team, error)
	UpdateTeam(ctx context.Context, updates UpdateTeam) (Team, error)
	CreateTeam(ctx context.Context, team CreateTeam) (Team, error)
}

type TeamBus struct {
	store TeamStorer
}

func NewTeamBus(store TeamStorer) *TeamBus {
	return &TeamBus{
		store: store,
	}
}

func (tb *TeamBus) GetTeamsBy(ctx context.Context, query QueryTeam) ([]Team, error) {
	teams, err := tb.store.GetTeamsBy(ctx, query)
	if err != nil {
		return []Team{}, fmt.Errorf("bus:GetTeamsBy:%w", err)
	}

	return teams, nil
}

func (tb *TeamBus) CreateTeam(ctx context.Context) (Team, error) {
	team, err := tb.store.CreateTeam(ctx, CreateTeam{
		Name:    "team" + rand.Text(),
		Country: "country" + rand.Text(),
	})

	if err != nil {
		return Team{}, fmt.Errorf("bus:CreateTeam:%w", err)
	}
	return team, nil

}

func (tb *TeamBus) UpdateTeam(ctx context.Context, updates UpdateTeam) (Team, error) {
	team, err := tb.store.UpdateTeam(ctx, updates)
	if err != nil {
		return Team{}, fmt.Errorf("bus:UpdateTeam:%w", err)
	}

	return team, nil
}
