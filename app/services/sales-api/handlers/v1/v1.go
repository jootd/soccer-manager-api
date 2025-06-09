// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/app/services/sales-api/handlers/v1/playergrp"
	"github.com/jootd/soccer-manager/app/services/sales-api/handlers/v1/teamgrp"
	"github.com/jootd/soccer-manager/app/services/sales-api/handlers/v1/transfergrp"
	"github.com/jootd/soccer-manager/app/services/sales-api/handlers/v1/usergrp"
	"github.com/jootd/soccer-manager/business/adapter/playeradapter"
	"github.com/jootd/soccer-manager/business/adapter/teamadapter"
	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/playerbus/stores/playerdb"
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/teambus/stores/teamdb"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
	"github.com/jootd/soccer-manager/business/domain/transferbus/stores/transferdb"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/domain/userbus/stores/userdb"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"github.com/jootd/soccer-manager/business/sdk/v1/mid"
	"github.com/jootd/soccer-manager/foundation/web"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *zap.SugaredLogger
	DB  *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	userStore := userdb.NewStore(cfg.Log, cfg.DB)
	teamStore := teamdb.NewStore(cfg.Log, cfg.DB)
	playerStore := playerdb.NewStore(cfg.Log, cfg.DB)
	transferStore := transferdb.NewStore(cfg.Log, cfg.DB)

	userBus := userbus.NewUserBus(userStore, cfg.Log, nil)
	teamBus := teambus.NewTeamBus(teamStore, cfg.Log, nil)

	playerAdapter := playeradapter.NewAdapter(playerStore)
	teamAdapter := teamadapter.NewAdapter(teamStore)

	playerBus := playerbus.NewPlayerBus(playerStore, cfg.Log, nil)
	transferBus := transferbus.NewTransferBus(transferStore, cfg.Log, playerAdapter, teamAdapter)

	ugh := usergrp.Handlers{
		UserBus:   userBus,
		TeamBus:   teamBus,
		PlayerBus: playerBus,
		Tx:        sqldb.NewBeginner(cfg.DB),
	}
	app.Handle(http.MethodPost, version, "/auth/signin", ugh.Signin)
	app.Handle(http.MethodPost, version, "/auth/signup", ugh.Signup)

	team := teamgrp.Handlers{
		UserBus:   userBus,
		TeamBus:   teamBus,
		PlayerBus: playerBus,
	}

	authMid := mid.Authorize(cfg.Log, userBus, teamBus)

	app.Handle(http.MethodGet, version, "/team", team.Get, authMid)
	app.Handle(http.MethodPatch, version, "/team", team.Update, authMid)

	player := playergrp.Handlers{
		TeamBus:   teamBus,
		PlayerBus: playerBus,
	}
	app.Handle(http.MethodGet, version, "/player", player.All, authMid)
	app.Handle(http.MethodGet, version, "/player/:id", player.ById, authMid)

	transfer := transfergrp.Handlers{
		PlayerBus:   playerBus,
		TransferBus: transferBus,
	}

	_ = transfer

}
