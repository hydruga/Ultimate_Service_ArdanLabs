package testgrp

import (
	"context"
	"encoding/json"
	"net/http"

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

	statusCode := http.StatusOK

	// Don't want to log inside of handlers.
	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
	return json.NewEncoder(w).Encode(status)
}
