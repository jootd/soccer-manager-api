package playerbus

import (
	"context"
	"crypto/rand"
	"fmt"
	mrand "math/rand"

	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"github.com/jootd/soccer-manager/business/types/age"
	"github.com/jootd/soccer-manager/business/types/position"
	"go.uber.org/zap"
)

const (
	InitialPlayerValue      = 1_000_000 //$
	InitialGoalkeepersCount = 3
	InitialDefendersCount   = 6
	InitialMidfieldersCount = 6
	InitialAttackers        = 5
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	All(ctx context.Context) ([]Player, error)
	GetByTeamID(ctx context.Context, teamID int) ([]Player, error)
	Query(ctx context.Context, query QueryFilter) ([]Player, error)
	Update(ctx context.Context, player Player) error
	Create(ctx context.Context, new Player) (int, error)
	CreateBatch(ctx context.Context, players []Player) error
}

type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, new CreatePlayer) (Player, error)
	GenerateInitialBatch(ctx context.Context, teamID int) error
	Update(ctx context.Context, player UpdatePlayer) error
	Query(ctx context.Context, filter QueryFilter) ([]Player, error)
	GetByTeamID(ctx context.Context, teamID int) ([]Player, error)
}

type Extension func(ExtBusiness) ExtBusiness

type Business struct {
	store Storer
	log   *zap.SugaredLogger
}

func NewPlayerBus(store Storer, log *zap.SugaredLogger, extensions ...Extension) ExtBusiness {
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

func (up *Business) GetByTeamID(ctx context.Context, teamID int) ([]Player, error) {
	result, err := up.store.GetByTeamID(ctx, teamID)
	if err != nil {
		return []Player{}, fmt.Errorf("playerbus:GetByTeamID:%w", err)
	}
	return result, nil
}

func (up *Business) Create(ctx context.Context, new CreatePlayer) (Player, error) {
	newPlayer := Player{
		TeamID:    new.TeamID,
		FirstName: new.FirstName,
		LastName:  new.LastName,
		Age:       new.Age,
		Country:   new.Country,
		Value:     new.Value,
		Position:  new.Position,
	}
	id, err := up.store.Create(ctx, newPlayer)

	if err != nil {
		return Player{}, fmt.Errorf("playerbus:Create:%w", err)
	}
	newPlayer.ID = id
	return newPlayer, nil
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

	var players []Player

	for _, rc := range roles {
		for i := 0; i < rc.count; i++ {
			min := 18
			max := 40

			randomAge := mrand.Intn(max-min+1) + min

			player := Player{
				TeamID:    teamID,
				FirstName: rand.Text(),
				LastName:  rand.Text(),
				Country:   rand.Text()[:2],
				Age:       age.MustParse(randomAge),
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

func (up *Business) Update(ctx context.Context, update UpdatePlayer) error {
	player := Player{
		ID: update.ID,
	}

	if update.TeamID != nil {
		player.TeamID = *update.TeamID
	}

	if update.Country != nil {
		player.Country = *update.Country
	}

	if update.FirstName != nil {
		player.FirstName = *update.FirstName
	}

	if update.LastName != nil {
		player.LastName = *update.LastName
	}

	if update.Age != nil {
		player.Age = *update.Age
	}

	if update.Value != nil {
		player.Value = *update.Value
	}

	err := up.store.Update(ctx, player)
	if err != nil {
		return fmt.Errorf("playerbus:Update:%w", err)
	}

	return nil
}

func (up *Business) Query(ctx context.Context, filter QueryFilter) ([]Player, error) {
	players, err := up.store.Query(ctx, filter)
	if err != nil {
		return []Player{}, fmt.Errorf("playerbus:Query:%w", err)
	}
	return players, nil
}
