package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := chi.URLParam(r, "*")
		if len(h.AllowedPrefixes) > 0 {
			ok := false
			for _, prefix := range h.AllowedPrefixes {
				if strings.HasPrefix(path, prefix) {
					ok = true
				}
			}
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
		url := h.PximgRoot + path
		contentType, body, err := h.Client.Download(r.Context(), url, false)
		if err != nil {
			h.reportError(w, err)
			return
		}
		w.Header().Add("Content-Type", contentType)
		w.Write(body)
	}
}
