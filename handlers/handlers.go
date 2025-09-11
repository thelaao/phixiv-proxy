package handlers

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/thelaao/phixiv-proxy/pixiv"
)

type Handler struct {
	Client            *pixiv.PixivClient
	AllowedPrefixes   []string
	UgoiraMaxDuration int
	UgoiraMinDuration int
	UgoiraMaxFrames   int
	PximgRoot         *url.URL
}

func NewHandler(client *pixiv.PixivClient, allowedPrefixes string, ugoiraMaxDuration int, ugoiraMinDuration int, ugoiraMaxFrames int, pximgRoot string) *Handler {
	if len(pximgRoot) == 0 {
		pximgRoot = "https://i.pximg.net/"
	}
	pxImgRootUrl, err := url.Parse(pximgRoot)
	if err != nil {
		log.Fatalf("invalid url %s", pximgRoot)
	}
	return &Handler{
		Client:            client,
		AllowedPrefixes:   strings.Split(allowedPrefixes, ","),
		UgoiraMaxDuration: ugoiraMaxDuration,
		UgoiraMinDuration: ugoiraMinDuration,
		UgoiraMaxFrames:   ugoiraMaxFrames,
		PximgRoot:         pxImgRootUrl,
	}
}

func (h *Handler) reportError(w http.ResponseWriter, err error) {
	returnCode := http.StatusInternalServerError
	if perr, ok := err.(*pixiv.PixivRequestError); ok {
		returnCode = perr.StatusCode
	}
	http.Error(w, err.Error(), returnCode)
}
