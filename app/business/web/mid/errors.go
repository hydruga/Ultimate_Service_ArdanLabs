package mid

import (
	"context"
	"net/http"

	"github.com/hydruga/ultimate_service/app/business/sys/validate"
	"github.com/hydruga/ultimate_service/app/foundation/web"
	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			// Run the next handler and catch any propogated error.
			if err := handler(ctx, w, r); err != nil {

				// Log the error
				log.Errorw("ERROR", "traceid", v.TraceID, "ERROR", err)

				// Build out error response.
				var er validate.ErrorResponse
				var status int
				switch act := validate.Cause(err).(type) {
				case *validate.RequestError:
					er = validate.ErrorResponse{
						Error: act.Error(),
					}
					status = act.Status
				case *validate.FieldErrors:
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: act.Error(),
					}
					status = http.StatusBadRequest
				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we recieve the shutdown err we need to return it
				// back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}
			return nil
		}
		return h
	}
	return m
}
