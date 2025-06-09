package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	v1Web "github.com/jootd/soccer-manager/business/sdk/v1"
	"github.com/jootd/soccer-manager/business/sdk/v1/jwt"
	"github.com/jootd/soccer-manager/foundation/web"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	UserBus   userbus.ExtBusiness
	TeamBus   teambus.ExtBusiness
	PlayerBus playerbus.ExtBusiness
	Tx        sqldb.Beginner
}

func (h Handlers) Signup(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var newUser userbus.CreateUser

	if err := web.Decode(r, &newUser); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	tx, err := h.Tx.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("unable to start tx: %w", err)
	}
	userBusTx, err := h.UserBus.NewWithTx(tx)
	if err != nil {
		return fmt.Errorf("userbus.NewWithTx failed: %w", err)
	}

	teamBusTx, err := h.TeamBus.NewWithTx(tx)
	if err != nil {
		return fmt.Errorf("teambus.NewWithTx failed: %w", err)
	}

	playerTx, err := h.PlayerBus.NewWithTx(tx)
	if err != nil {
		return fmt.Errorf("teambus.NewWithTx failed: %w", err)
	}

	err = userBusTx.Create(ctx, newUser)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueUsername) {
			return v1Web.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("user[%+v]: %w", &newUser, err)
	}

	newTeam, err := teamBusTx.AutoGenerate(ctx)
	if err != nil {
		return fmt.Errorf("handler:Signup:%w", err)
	}

	err = playerTx.GenerateInitialBatch(ctx, newTeam.ID)
	if err != nil {
		return fmt.Errorf("handler.Signup:%w", err)
	}

	err = userBusTx.Update(ctx, userbus.UpdateUser{Username: &newUser.Username, TeamID: &newTeam.ID})
	if err != nil {
		return fmt.Errorf("handler.Signup:%w", err)
	}

	account := struct {
		User userbus.CreateUser `json:"user"`
		Team teambus.Team       `json:"team"`
	}{
		User: newUser,
		Team: newTeam,
	}

	tx.Commit()
	return web.Respond(ctx, w, account, http.StatusCreated)

}

func (h Handlers) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var login userbus.CreateUser

	if err := web.Decode(r, &login); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	user, err := h.UserBus.Get(r.Context(), login.Username)
	if err != nil {
		if errors.Is(err, userbus.ErrNotFound) {
			return v1Web.NewRequestError(err, http.StatusNotFound)
		}
	}

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		return v1Web.NewRequestError(err, http.StatusUnauthorized)
	}

	token, err := jwt.GenerateJWT(user.Username)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("signin: usr[%+v]: %s", login, err.Error()), http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "donottouchme",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		// Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	return web.Respond(ctx, w, "OK", http.StatusOK)
}
