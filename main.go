package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/isaac-albert/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = ":8080"

type apiConfig struct {
	fileServerhits atomic.Int32
	dbQuery *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error opening a connection to Data base")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	mux := http.NewServeMux()

	//filepathRoot := http.Dir(".")
	assetpathRoot := http.Dir(".")
	hand := http.StripPrefix("/app", http.FileServer(assetpathRoot))
	apiCfg := &apiConfig{
		fileServerhits: atomic.Int32{},
		dbQuery: dbQueries,
	}
	apiCfg.fileServerhits.Store(0)

	//GetMiddlware(apiCfg.middlewareMetricsInc(hand))
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(hand))
	mux.HandleFunc("GET /api/healthz", http.HandlerFunc(handlerReadiness))
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(apiCfg.handlerReset))
	mux.HandleFunc("POST /api/validate_chirp", http.HandlerFunc(ParseJson))

	// /http.FileServer(assetpathRoot)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", assetpathRoot, port)
	log.Fatal(srv.ListenAndServe())
}