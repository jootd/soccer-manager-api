package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jootd/soccer-manager/business/domain/userbus"
	v1Web "github.com/jootd/soccer-manager/business/sdk/v1"
	"github.com/jootd/soccer-manager/business/sdk/v1/jwt"
	"github.com/jootd/soccer-manager/foundation/web"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	Bus userbus.ExtBusiness
}

func (h Handlers) Signup(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var newUser userbus.CreateUser

	if err := web.Decode(r, &newUser); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	err := h.Bus.Create(ctx, newUser)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueUsername) {
			return v1Web.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("user[%+v]: %w", &newUser, err)
	}

	return web.Respond(ctx, w, newUser, http.StatusCreated)

}

func (h Handlers) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var login userbus.CreateUser

	if err := web.Decode(r, &login); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	user, err := h.Bus.Get(r.Context(), login.Username)
	if err != nil {
		if errors.Is(err, userbus.ErrNotFound) {
			return v1Web.NewRequestError(err, http.StatusNotFound)
		}
	}

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		return v1Web.NewRequestError(err, http.StatusUnauthorized)
	}

	token, err := jwt.GenerateJWT(user.Password)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("signin: usr[%+v]: %s", login, err.Error()), http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "donottouchme",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	return web.Respond(ctx, w, "OK", http.StatusOK)
}
