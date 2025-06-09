package transfergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
	v1Web "github.com/jootd/soccer-manager/business/sdk/v1"
	"github.com/jootd/soccer-manager/business/sdk/v1/mid"
	"github.com/jootd/soccer-manager/foundation/web"
)

var (
	ErrInvalidID        = errors.New("ID is not in its proper form")
	ErrResourceNotFound = errors.New("resource not found")
)

type Handlers struct {
	PlayerBus   playerbus.ExtBusiness
	TransferBus transferbus.ExtBusiness
}

func (h Handlers) All(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	players, err := h.TransferBus.All(ctx)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusNotFound)
	}

	return web.Respond(ctx, w, players, http.StatusOK)
}

func (h Handlers) Buy(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	transferId, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return v1Web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	val := ctx.Value(mid.TeamIdContextKey)
	buyerTeamID := val.(int)
	err = h.TransferBus.Buy(ctx, transferId, buyerTeamID)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusConflict)
	}

	return web.Respond(ctx, w, "OK", http.StatusOK)

}

func (h Handlers) ListForSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	price := struct {
		AskingPrice int64 `json:"asking_price"`
	}{}

	if err := web.Decode(r, &price); err != nil {
		nrr := fmt.Errorf("unable to decode payload: %w", err)
		return web.Respond(ctx, w, nrr, http.StatusBadRequest)
	}

	playerID, err := strconv.Atoi(web.Param(r, "player"))
	if err != nil {
		return v1Web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	val := ctx.Value(mid.TeamIdContextKey)
	sellerID := val.(int)

	err = h.TransferBus.ListForSale(ctx, playerID, sellerID, price.AskingPrice)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusBadRequest)
	}

	return web.Respond(ctx, w, "OK", http.StatusOK)

}

func (h Handlers) ById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	playerId, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return v1Web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	result, err := h.PlayerBus.Query(ctx, playerbus.QueryFilter{
		ID: &playerId,
	})
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusNotFound)
	}
	return web.Respond(ctx, w, result[0], http.StatusOK)

}

func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	playerId, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return v1Web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}
	var upd playerbus.UpdatePlayer
	if err := web.Decode(r, &upd); err != nil {
		nrr := fmt.Errorf("unable to decode payload: %w", err)
		return web.Respond(ctx, w, nrr, http.StatusBadRequest)
	}

	upd.ID = playerId
	err = h.PlayerBus.Update(r.Context(), upd)
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusInternalServerError)
	}

	result, err := h.PlayerBus.Query(ctx, playerbus.QueryFilter{ID: &playerId})
	if err != nil {
		return web.Respond(ctx, w, err, http.StatusNotFound)
	}

	return web.Respond(ctx, w, result[0], http.StatusOK)
}
