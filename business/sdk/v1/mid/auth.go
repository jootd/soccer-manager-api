package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/sdk/v1/jwt"
	"github.com/jootd/soccer-manager/foundation/web"
	"go.uber.org/zap"
)

type ContextKey string

const (
	UsernameContextKey ContextKey = "ctxusername"
	TeamIdContextKey   ContextKey = "ctxteamid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrKeyMissing   = errors.New("key is missing")
)

func Authorize(log *zap.SugaredLogger, userBus userbus.ExtBusiness, teamBus teambus.ExtBusiness) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			cookie, err := r.Cookie("donottouchme")
			if err != nil {
				return errors.New("unauthorized: missing or invalid cookie")
			}

			tokenString := cookie.Value
			username, err := jwt.ValidateJWT(tokenString)
			if err != nil {
				// http.Error(w, "unauthorized", http.StatusUnauthorized)
				return errors.Join(err, ErrUnauthorized)
			}

			fmt.Println("user name is ", username)

			user, err := userBus.Get(r.Context(), username)
			if err != nil {
				// http.Error(w, "unauthorized", http.StatusUnauthorized)
				return errors.Join(err, ErrUnauthorized)
			}

			fmt.Printf("========user is %+v\n", user)

			ctx = context.WithValue(ctx, UsernameContextKey, username)
			ctx = context.WithValue(ctx, TeamIdContextKey, user.TeamID)

			fmt.Printf("======>><<<====new context is %+v\n", ctx)
			r = r.WithContext(ctx)

			return handler(ctx, w, r)

		}
		return h
	}

	return m
}
