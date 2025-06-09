package transferbus

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"github.com/jootd/soccer-manager/business/types/transferstatus"
	"go.uber.org/zap"
)

const (
	MinFactor = 0.1
	MaxFactor = 1
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	All(ctx context.Context) ([]Transfer, error)
	Create(ctx context.Context, t Transfer) (int, error)
	Query(ctx context.Context, query QueryFilter) ([]Transfer, error)
	Update(ctx context.Context, update Transfer) error
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

type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	All(ctx context.Context) ([]Transfer, error)
	ListForSale(ctx context.Context, playerID, sellerID int, askingPrice int64) error
	Buy(ctx context.Context, transferID, buyerID int) error
}

type Extension func(ExtBusiness) ExtBusiness

type Business struct {
	store         Storer
	log           *zap.SugaredLogger
	playerService PlayerService
	teamService   TeamService
}

func NewTransferBus(store Storer, log *zap.SugaredLogger, ps PlayerService, ts TeamService, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		store:         store,
		log:           log,
		playerService: ps,
		teamService:   ts})

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

func (tc *Business) All(ctx context.Context) ([]Transfer, error) {
	transfers, err := tc.store.All(ctx)
	if err != nil {
		return []Transfer{}, fmt.Errorf("transferbus:market:%w", err)
	}
	return transfers, nil
}

func (tc *Business) ListForSale(ctx context.Context, playerID, sellerID int, askingPrice int64) error {
	//TODO: validation
	_, err := tc.store.Create(ctx, Transfer{
		PlayerID:    playerID,
		SellerID:    sellerID,
		AskingPrice: askingPrice,
		Status:      transferstatus.Listed,
	})

	if err != nil {
		return fmt.Errorf("transferbus:ListForSale:%w", err)
	}
	return nil
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
	err := tc.store.Update(ctx, Transfer{
		ID:     transferID,
		Status: transferstatus.Sold,
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
