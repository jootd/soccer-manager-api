package business

import "context"

type PlayerStorer interface {
	GetPlayersBy(ctx context.Context, query QueryPlayer) ([]Player, bool)
	UpdatePlayer(ctx context.Context, player UpdatePlayer) (Player, bool)
}

type PlayerBusiness struct {
}

type Player struct {
}

type QueryPlayer struct {
	teamId *int
}

type UpdatePlayer struct {
}
