package api

import (
	"fmt"
	"net/http"
)


type ApiConfig struct {
    fileserverHits int
}


func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

        cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) HitHandler( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    hitCount := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
	w.Write([]byte(hitCount))
}

func (cfg *ApiConfig) ResetHandler ( w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

