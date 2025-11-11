package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const port = ":8080"

type apiConfig struct {
	fileServerhits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	//filepathRoot := http.Dir(".")
	assetpathRoot := http.Dir(".")
	hand := http.StripPrefix("/app", http.FileServer(assetpathRoot))
	apiCfg := &apiConfig{
		fileServerhits: atomic.Int32{},
	}
	apiCfg.fileServerhits.Store(0)

	//GetMiddlware(apiCfg.middlewareMetricsInc(hand))
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(hand))
	mux.HandleFunc("GET /api/healthz", http.HandlerFunc(handlerReadiness))
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(apiCfg.handlerReset))

	// /http.FileServer(assetpathRoot)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", assetpathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	htmlSnippet := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerhits.Load())
	w.Write([]byte(htmlSnippet))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerhits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// ...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerhits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
