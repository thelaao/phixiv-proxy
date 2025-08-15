package handlers

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := h.PximgRoot
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		_, _, err := h.Client.Download(ctx, url, true)
		if err != nil {
			h.reportError(w, err)
			return
		}
		w.Header().Add("X-Pixiv-Target", url)
		w.WriteHeader(http.StatusOK)
	}
}
