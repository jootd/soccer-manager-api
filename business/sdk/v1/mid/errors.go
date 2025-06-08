package mid

import (
	"context"
	"net/http"

	v1Web "github.com/jootd/soccer-manager/business/sdk/v1"
	"github.com/jootd/soccer-manager/foundation/web"
	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				ctx, span := web.AddSpan(ctx, "business.web.v1.mid.error")
				span.RecordError(err)
				span.End()

				var er v1Web.ErrorResponse
				var status int
				switch {
				case v1Web.IsRequestError(err):
					reqErr := v1Web.GetRequestError(err)
					er = v1Web.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				// case auth.IsAuthError(err):
				// 	er = v1Web.ErrorResponse{
				// 		Error: http.StatusText(http.StatusUnauthorized),
				// 	}
				// 	status = http.StatusUnauthorized

				default:
					er = v1Web.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
