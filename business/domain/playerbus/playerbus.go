package playerbus

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/jootd/soccer-manager/business/types/age"
	"github.com/jootd/soccer-manager/business/types/position"
)

const (
	InitialPlayerValue      = 1_000_000 //$
	InitialGoalkeepersCount = 3
	InitialDefendersCount   = 6
	InitialMidfieldersCount = 6
	InitialAttackers        = 5
)

type Storer interface {
	Query(ctx context.Context, query QueryFilter) ([]Player, error)
	Update(ctx context.Context, player UpdatePlayer) (Player, error)
	Create(ctx context.Context, new CreatePlayer) (Player, error)
	CreateBatch(ctx context.Context, players []CreatePlayer) error
}

type Business struct {
	store Storer
}

func NewPlayerBus(store Storer) *Business {
	return &Business{
		store: store,
	}
}

func (up *Business) Create(ctx context.Context, new CreatePlayer) (Player, error) {
	player, err := up.store.Create(ctx, new)
	if err != nil {
		return Player{}, fmt.Errorf("playerbus:Create:%w", err)
	}
	return player, nil
}

func (up *Business) CreateBatch(ctx context.Context, players []CreatePlayer) error {
	if len(players) == 0 {
		return nil
	}
	if err := up.store.CreateBatch(ctx, players); err != nil {
		return fmt.Errorf("playerbus:CreateBatch:%w", err)
	}
	return nil
}

func (up *Business) GenerateInitialBatch(ctx context.Context, teamID int) error {
	roles := []struct {
		position position.Position
		count    int
	}{
		{position.Goalkeeper, 3},
		{position.Defender, 6},
		{position.Midfielder, 6},
		{position.Attacker, 5},
	}

	var players []CreatePlayer

	for _, rc := range roles {
		for i := 0; i < rc.count; i++ {
			player := CreatePlayer{
				TeamID:    teamID,
				FirstName: rand.Text(),
				LastName:  rand.Text(),
				Country:   rand.Text()[:2],
				Age:       age.MustParse(18),
				Value:     InitialPlayerValue,
				Position:  rc.position,
			}
			players = append(players, player)
		}
	}

	if err := up.store.CreateBatch(ctx, players); err != nil {
		return fmt.Errorf("playerbus: GenerateInitialBatch: %w", err)
	}

	return nil

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
