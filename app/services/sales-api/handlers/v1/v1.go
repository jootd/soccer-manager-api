// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/app/services/sales-api/handlers/v1/usergrp"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/domain/userbus/stores/userdb"
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
	ugh := usergrp.Handlers{
		Bus: userbus.NewUserBus(userdb.NewStore(cfg.Log, cfg.DB), cfg.Log, nil),
	}
	app.Handle(http.MethodPost, version, "/auth/signin", ugh.Signin)
	app.Handle(http.MethodPost, version, "/auth/signup", ugh.Signup)

}
