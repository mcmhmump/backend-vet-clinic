package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mcmhmump/backend-vet-clinic/internal/usecase"
)

type CacheHandler struct {
	cache *usecase.CacheService
}

func NewCacheHandler(cache *usecase.CacheService) *CacheHandler {
	return &CacheHandler{cache: cache}
}

func (h *CacheHandler) InvalidateCache(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	w.Header().Set("Content-Type", "application/json")

	if key == "" {
		h.cache.ClearAll()
		json.NewEncoder(w).Encode(map[string]string{
			"status": "all_cache_cleared",
		})
		return
	}

	h.cache.Invalidate(key)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "key_invalidated",
		"key":    key,
	})
}

type ResponseRecorder struct {
	http.ResponseWriter
	Body   []byte
	Status int
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body = append(r.Body, b...)
	return r.ResponseWriter.Write(b)
}

func CacheMiddleware(cache *usecase.CacheService, ttl time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		key := r.URL.Path + "?" + r.URL.RawQuery

		if cachedData, found := cache.Get(key); found {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.Write(cachedData)
			return
		}

		recorder := &ResponseRecorder{
			ResponseWriter: w,
			Body:           []byte{},
			Status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		if recorder.Status >= 200 && recorder.Status < 300 {
			cache.Set(key, recorder.Body, ttl)
		}
	}
}
