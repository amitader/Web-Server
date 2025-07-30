package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"database/sql"
	"sync/atomic"
	"github.com/amitader/web-Server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db  *database.Queries
	fileserverHits atomic.Int32
	platform string
	jwtSecret string
	polkaKey string
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	plat := os.Getenv("PLATFORM")
	if plat == "" {
		log.Fatal("PLATFORM must be set")
	}
	SECRET := os.Getenv("SECRET")
	if SECRET == "" {
		log.Fatal("SECRET must be set")
	}

	POLKA_KEY := os.Getenv("POLKA_KEY")
	if POLKA_KEY == "" {
		log.Fatal("POLKA_KEY must be set")
	}
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform: plat,
		jwtSecret: SECRET,
		polkaKey: POLKA_KEY,
	}
	
	
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreation)
	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpsCreation)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.ChirpsDeletion)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.revoke)
	mux.HandleFunc("PUT /api/users", apiCfg.changeUserDetails)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.webhooks)
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(html))
}



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


