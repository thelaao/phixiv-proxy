package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	RedisClient *redis.Client
	Enabled     bool
}

type CachedResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func NewCache(redisUrl string) *Cache {
	if len(redisUrl) == 0 {
		return &Cache{
			Enabled: false,
		}
	}
	return &Cache{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: redisUrl,
		}),
		Enabled: true,
	}
}

func (cache *Cache) Query(ctx context.Context, key string) []byte {
	if !cache.Enabled {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	val, err := cache.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	return []byte(val)
}

func (cache *Cache) Save(ctx context.Context, key string, value []byte, expires bool) {
	if !cache.Enabled {
		return
	}
	var ttl time.Duration
	if expires {
		ttl = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	cache.RedisClient.Set(ctx, key, string(value), ttl)
}

func (cache *Cache) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/i/") {
			next.ServeHTTP(w, r)
			return
		}
		cached := cache.Query(r.Context(), r.URL.Path)
		if len(cached) > 0 {
			resp, err := FromBytes(cached)
			if err == nil {
				resp.CopyToWriter(w, r)
				return
			}
		}
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)
		resp := FromRecorder(recorder)
		resp.CopyToWriter(w, r)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
			resp.Save(r.Context(), cache, r.URL.Path, resp.StatusCode != http.StatusOK)
		}
	}
	return http.HandlerFunc(fn)
}

func (resp *CachedResponse) CopyToWriter(w http.ResponseWriter, r *http.Request) {
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	w.Header().Set("Accept-Ranges", "bytes")
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		w.Write(resp.Body)
		return
	}
	content := bytes.NewReader(resp.Body)
	http.ServeContent(w, r, "", time.Time{}, content)
}

func (resp *CachedResponse) Save(ctx context.Context, cache *Cache, path string, expires bool) {
	buffer := &bytes.Buffer{}
	gob.NewEncoder(buffer).Encode(resp)
	cache.Save(ctx, path, buffer.Bytes(), expires)
}

func FromRecorder(recorder *httptest.ResponseRecorder) *CachedResponse {
	return &CachedResponse{
		StatusCode: recorder.Code,
		Header:     recorder.Header(),
		Body:       recorder.Body.Bytes(),
	}
}

func FromBytes(data []byte) (resp *CachedResponse, err error) {
	resp = &CachedResponse{}
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(resp)
	return
}
