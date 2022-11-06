package testgrp

import (
	"context"
	"net/http"

	"github.com/hydruga/ultimate_service/app/foundation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "Ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
