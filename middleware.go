package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/userbus"
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

func CreateAuthMiddleware(userbus *userbus.Business, teamBus *teambus.Business) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("donottouchme")
			if err != nil {
				http.Error(w, "unauthorized: missing or invalid cookie", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			username, err := validateJWT(tokenString)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userbus.Get(r.Context(), username)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			teams, err := teamBus.Query(r.Context(), teambus.QueryFilter{ID: &user.TeamID})
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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
