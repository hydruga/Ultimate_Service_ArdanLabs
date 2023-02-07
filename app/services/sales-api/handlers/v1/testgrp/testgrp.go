// Package testgrp contains all the test handlers.
package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/hydruga/ultimate_service/app/business/sys/validate"
	"github.com/hydruga/ultimate_service/app/foundation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		//return errors.New("untrusted errors")
		//return web.NewShutdownError("restart service")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		//panic("testing panic")
	}

	status := struct {
		Status string
	}{
		Status: "Ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
