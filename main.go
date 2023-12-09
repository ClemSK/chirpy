package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ClemSK/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

const dbFilePath = "database.json"

func main() {
	const port = "8080"
	const filepathRoot = "."

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		// If --debug flag is provided, delete the database file
		err := os.Remove(dbFilePath)
		if err != nil {
			fmt.Println("Error deleting database file:", err)
			return
		}
		fmt.Println("Database file deleted (debug mode).")
	}

	// Check if the database file exists
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		// Create an empty JSON object and write it to the file
		data := make(map[string]interface{})
		err := writeJsonFile(dbFilePath, data)
		if err != nil {
			fmt.Println("Error creating database file:", err)
			return
		}
		fmt.Println("Database file created successfully.")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	r := chi.NewRouter() // base router
	r.Use(middleware.Logger)

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler) // satisfying the chi router for routes w/out trailing slash
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter() // api router
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handleReset)

	apiRouter.Post("/login", apiCfg.handlerLogin)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Get("/users", apiCfg.handlerUserGet)
	apiRouter.Get("/users/{id}", apiCfg.handlerUserGetById)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsGet)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerChirpsGetById)
	r.Mount("/api", apiRouter) // using the sub-router

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", adminRouter) // using the sub-router

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
