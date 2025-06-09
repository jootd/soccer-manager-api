package teamadapter

import (
	"context"
	"fmt"

	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
)

type Storer interface {
	Query(ctx context.Context, filter teambus.QueryFilter) ([]teambus.Team, error)
	Update(ctx context.Context, updates teambus.Team) error
	GetByID(ctx context.Context, id int) (teambus.Team, error)
}

type Adapter struct {
	store Storer
}

func NewAdapter(store Storer) *Adapter {
	return &Adapter{
		store: store,
	}
}

func (t *Adapter) GeTeamInfo(ctx context.Context, id int) (transferbus.TeamInfo, error) {
	teams, err := t.store.Query(ctx, teambus.QueryFilter{ID: &id})
	if err != nil {
		return transferbus.TeamInfo{}, fmt.Errorf("teamadapter:GetTeamInfo:%w", err)
	}

	return toTransferTeamInfo(teams[0]), nil
}

func (t *Adapter) UpdateBudget(ctx context.Context, teamID int, newBudget int64) error {
	err := t.store.Update(ctx, teambus.Team{
		ID:     teamID,
		Budget: newBudget})
	if err != nil {
		return fmt.Errorf("teamadapter:UpdateBudget:%w", err)
	}

	return nil
}
func (t *Adapter) GetByID(ctx context.Context, id int) (transferbus.TeamInfo, error) {
	//TODO:
	return transferbus.TeamInfo{}, nil

}
