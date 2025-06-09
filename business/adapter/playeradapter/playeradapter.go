package playeradapter

import (
	"context"
	"fmt"

	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
)

type Storer interface {
	Query(ctx context.Context, query playerbus.QueryFilter) ([]playerbus.Player, error)
	Update(ctx context.Context, player playerbus.Player) error
}

type Adapter struct {
	store Storer
}

func NewAdapter(store Storer) *Adapter {
	return &Adapter{
		store: store,
	}
}

func (a *Adapter) GetPlayerInfo(ctx context.Context, id int) (transferbus.PlayerInfo, error) {
	players, err := a.store.Query(ctx, playerbus.QueryFilter{ID: &id})
	if err != nil {
		return transferbus.PlayerInfo{}, fmt.Errorf("playeradapter:GetPlayerInfo:%w", err)
	}

	return toTransferInfo(players[0]), nil
}

// updates team for player
func (a *Adapter) UpdateTeam(ctx context.Context, playerID, teamID int) error {
	err := a.store.Update(ctx, playerbus.Player{
		ID:     playerID,
		TeamID: teamID,
	})

	if err != nil {
		return fmt.Errorf("playeradapter:UpdateTeam:%w", err)
	}

	return nil
}

func (a *Adapter) UpdateValue(ctx context.Context, playerID int, newValue int64) error {
	err := a.store.Update(ctx, playerbus.Player{
		ID:    playerID,
		Value: newValue,
	})
	if err != nil {
		return fmt.Errorf("playeradapter:UpdateValue:%w", err)
	}

	return nil
}
