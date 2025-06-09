package teamgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/sdk/v1/mid"
	"github.com/jootd/soccer-manager/business/view/vteambus"
	"github.com/jootd/soccer-manager/foundation/web"
)

type Handlers struct {
	UserBus   userbus.ExtBusiness
	TeamBus   teambus.ExtBusiness
	PlayerBus playerbus.ExtBusiness
}

func (h Handlers) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	val := ctx.Value(mid.TeamIdContextKey)

	teamId := val.(int)
	players, err := h.PlayerBus.GetByTeamID(ctx, teamId)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusNotFound)
	}

	team, err := h.TeamBus.GetByID(ctx, teamId)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusNotFound)
	}

	teamWithPlayers := vteambus.FromTeam(team, players)

	return web.Respond(ctx, w, teamWithPlayers, http.StatusOK)

}

func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var upd teambus.UpdateTeam

	val := ctx.Value(mid.TeamIdContextKey)

	teamId := val.(int)
	if err := web.Decode(r, &upd); err != nil {
		nrr := fmt.Errorf("unable to decode payload: %w", err)
		return web.Respond(ctx, w, nrr, http.StatusBadRequest)
	}

	upd.ID = teamId
	err := h.TeamBus.Update(r.Context(), upd)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusInternalServerError)
	}

	team, err := h.TeamBus.GetByID(ctx, teamId)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusNotFound)
	}

	return web.Respond(ctx, w, team, http.StatusOK)
}
