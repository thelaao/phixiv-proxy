package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thelaao/phixiv-proxy/handlers"
	"github.com/thelaao/phixiv-proxy/pixiv"
	"github.com/thelaao/phixiv-proxy/utils"
)

func main() {
	client := pixiv.NewPixivClient(&http.Client{}, os.Getenv("USER_AGENT"), os.Getenv("PIXIV_COOKIE"))
	cache := utils.NewCache(os.Getenv("REDIS_URL"))
	h := handlers.NewHandler(client, os.Getenv("ALLOWED_PREFIXES"), parseEnvInt("UGOIRA_DURATION"), parseEnvInt("UGOIRA_MIN_DURATION"), parseEnvInt("UGOIRA_FRAMES"), os.Getenv("PXIMG_BASE"))

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cache.Middleware)

	r.Get(`/i/ugoira/{id:\d+}.{format:[a-z0-9]+}`, h.UgoiraHandler())
	r.Get(`/i/*`, h.ProxyHandler())
	r.Get(`/health`, h.HealthHandler())

	log.Fatal(http.ListenAndServe(":3000", r))
}

func parseEnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return value
}
