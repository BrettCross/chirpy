package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/brettcross/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	filepathRoot := "."
	port := "8080"

	// load environment variables
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	// open raw db connection
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening db connection: %s", err)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// register custom handlers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

	
}
