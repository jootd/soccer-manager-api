package transferbus

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jootd/soccer-manager/business/types/transferstatus"
)

const (
	MinFactor = 0.1
	MaxFactor = 1
)

type Storer interface {
	Create(ctx context.Context, t Transfer) error
	Query(ctx context.Context, query QueryFilter) ([]Transfer, error)
	Update(ctx context.Context, update UpdateTransfer) (Transfer, error)
	// Delete(ctx context.Context, id int) error
}

// we're gonna use playeradapter as PlayerService
type PlayerService interface {
	GetPlayerInfo(ctx context.Context, id int) (PlayerInfo, error)
	UpdateTeam(ctx context.Context, playerID, teamID int) error
	UpdateValue(ctx context.Context, playerID int, newValue int64) error
}

// same here
type TeamService interface {
	GetByID(ctx context.Context, id int) (TeamInfo, error)
	UpdateBudget(ctx context.Context, teamID int, newBudget int64) error
}

type Business struct {
	store         Storer
	playerService PlayerService
	teamService   TeamService
}

func New(store Storer, ps PlayerService, ts TeamService) *Business {
	return &Business{store, ps, ts}
}

func (tc *Business) All(ctx context.Context) ([]Transfer, error) {
	transfers, err := tc.store.Query(ctx, QueryFilter{})
	if err != nil {
		return []Transfer{}, fmt.Errorf("transferbus:market:%w", err)
	}
	return transfers, nil
}

func (tc *Business) ListForSale(ctx context.Context, playerID, sellerID int, askingPrice int64) error {
	//TODO: validation
	return tc.store.Create(ctx, Transfer{
		PlayerID:    playerID,
		SellerID:    sellerID,
		AskingPrice: askingPrice,
		Status:      transferstatus.Listed,
	})
}

func (tc *Business) Buy(ctx context.Context, transferID, buyerID int) error {
	transfers, err := tc.store.Query(ctx, QueryFilter{ID: &transferID})
	if err != nil {
		return err
	}
	transfer := transfers[0]

	// if not listed
	if !transfer.Status.Equal(transferstatus.Listed) {
		return errors.New("transferbus:Buy:player_not_available")
	}

	buyerTeam, err := tc.teamService.GetByID(ctx, buyerID)
	if err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	if buyerTeam.Budget < transfer.AskingPrice {
		return errors.New("transferbus:Buy:insufficient_budget")
	}

	// Transfer ownership
	if err := tc.playerService.UpdateTeam(ctx, transfer.PlayerID, buyerID); err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	// Adjust budgets
	sellerTeam, err := tc.teamService.GetByID(ctx, transfer.SellerID)
	if err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	// update buyer team budget
	if err := tc.teamService.UpdateBudget(ctx, buyerID, buyerTeam.Budget-transfer.AskingPrice); err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	// update seller team budget
	if err := tc.teamService.UpdateBudget(ctx, transfer.SellerID, sellerTeam.Budget+transfer.AskingPrice); err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	// caulcate player value
	newValue, err := tc.calculate(ctx, transfer.PlayerID)
	if err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	if err := tc.playerService.UpdateValue(ctx, transfer.PlayerID, newValue); err != nil {
		return fmt.Errorf("transferbus:Buy:%w", err)
	}

	// Mark as sold
	return tc.markAsSold(ctx, transferID)
}

func (tc *Business) markAsSold(ctx context.Context, transferID int) error {
	_, err := tc.store.Update(ctx, UpdateTransfer{
		ID:     transferID,
		Status: &transferstatus.Sold,
	})

	if err != nil {
		return fmt.Errorf("transferbus:MarkAsSold:%w", err)
	}
	return nil
}

// Increase by rand factor
func (tc *Business) calculate(ctx context.Context, playerID int) (int64, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	factor := r.Float64()*(MaxFactor-MinFactor) + MinFactor
	pinfo, err := tc.playerService.GetPlayerInfo(ctx, playerID)
	if err != nil {
		return 0, fmt.Errorf("transferbus:calculate:%w", err)
	}
	newValue := int64(float64(pinfo.Value) * (MaxFactor + factor))

	return newValue, nil
}
