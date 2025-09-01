package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/thelaao/phixiv-proxy/utils"
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
		if strings.HasPrefix(path, "img-master/img/") && contentType == "image/jpeg" {
			body = utils.ReencodeJPEG(body)
		}
		w.Header().Add("Content-Type", contentType)
		w.Header().Add("Content-Length", strconv.Itoa(len(body)))
		w.Write(body)
	}
}
