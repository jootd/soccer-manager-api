package business

import "context"

type Team struct {
	ID      int
	Name    string
	Country string
}

type UpdateTeam struct {
	Name    *string
	Country *string
}

type QueryTeam struct {
	ID      *int
	Name    *string
	Country *string
}

type TeamStorer interface {
	GetTeamsBy(ctx context.Context, query QueryTeam) ([]Team, bool)
	UpdateTeam(ctx context.Context, team UpdateTeam) (Team, bool)
	CreateTeam(ctx context.Context) (Team, bool)
}

type TeamBus struct {
	store TeamStorer
}

func NewTeamBus(store TeamStorer) *TeamBus {
	return &TeamBus{
		store: store,
	}
}

func (tb *TeamBus) GetTeamBy(ctx context.Context, query QueryTeam) (Team, error) {
	return Team{}, nil
}

func (tb *TeamBus) CreateTeam(ctx context.Context) (Team, error) {
	return Team{}, nil

}

func (tb *TeamBus) UpdateTeam(ctx context.Context, updateTeam UpdateTeam) (Team, error) {
	return Team{}, nil
}
