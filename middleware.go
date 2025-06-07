package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jootd/soccer-manager/business"
)

type ContextKey string

var (
	UsernameContextKey ContextKey = "ctxuname"
	UserContextKey     ContextKey = "ctxuser"
	UserTeamContextKey ContextKey = "ctxteam"
)

func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func CreateJWTMiddleware(userbus *business.UserBus, teambus *business.TeamBus) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Remove "Bearer " prefix
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			username, err := validateJWT(tokenString)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userbus.GetUser(r.Context(), username)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			teams, err := teambus.GetTeamsBy(r.Context(), business.QueryTeam{ID: &user.TeamId})
			if err != nil {
				http.Error(w, "team_not_found", http.StatusNotFound)
				return
			}

			team := teams[0]

			ctx := context.WithValue(r.Context(), UsernameContextKey, username)
			ctx = context.WithValue(ctx, UserContextKey, user)
			ctx = context.WithValue(ctx, UserTeamContextKey, team)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}
