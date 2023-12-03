package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Create a new ServeMux
	// mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler) // satisfying the chi router for routes w/out trailing slash
	r.Handle("/app/*", fsHandler)
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetrics)
	r.Get("/reset", apiCfg.handleReset)

	// Wrap the mux in the CORS middleware
	corsMux := middlewareCors(r)

	// Create a new HTTP server using the chi as the handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
