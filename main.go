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
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if dbg != nil && *dbg {
		// If --debug flag is provided, delete the database file
		err := db.ResetDB()
		if err != nil {
			fmt.Println("Error deleting database file:", err)
			return
		}
		fmt.Println("Database file deleted (debug mode).")
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	r := chi.NewRouter() // base router
	r.Use(middleware.Logger)

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler) // satisfying the chi router for routes w/out trailing slash
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter() // api router
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handleReset)

	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)
	apiRouter.Post("/login", apiCfg.handlerLogin)

	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Get("/users", apiCfg.handlerUserGet)
	apiRouter.Get("/users/{id}", apiCfg.handlerUserGetById)
	apiRouter.Put("/users", apiCfg.handlerUsersUpdate)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsGet)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerChirpsGetById)
	apiRouter.Delete("/chirps/{id}", apiCfg.handlerChirpsDelete)

	apiRouter.Post("/polka/webhooks", apiCfg.handleWebhook)

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
