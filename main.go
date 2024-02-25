package main

import (
	"log"
	"net/http"
	"os"

	// Using dot convention to becasue of clear Handler names
	"github.com/SisyphianLiger/Chirpy/api"
	"github.com/go-chi/chi/v5"
)


func main() {
    
    const filepathRoot = "./app"
    const port = "8080"

    fs := http.FileServer(http.Dir(filepathRoot))
    os.Remove("database.json")
    
    // Creating a empty struct to track state of Handler
    apiCfg := &api.ApiConfig{}

    // Creating a new App Router with chi
    r := chi.NewRouter()
    fsHandler := apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", fs))
    r.Handle("/app/*", fsHandler)
    r.Handle("/app", fsHandler)
    
    // Creating new API router with chi
    apiR := chi.NewRouter()
    // Ensuring only Get requests for /metrics and what not
    apiR.Get("/reset", apiCfg.ResetHandler)
    apiR.Get("/healthz", api.HealthzHandler)

    // REMOVE LATER: Just for testing

    dbrequester,_ := api.NewDB("database.json")
    apiR.Post("/chirps", dbrequester.HandlePostChirp)
    apiR.Get("/chirps", dbrequester.HandleGetAllChirps)
    apiR.Get("/chirps/{chirpID}", dbrequester.HandleGetID)

    
    // Mounting API to router
    r.Mount("/api", apiR)

    // Mounting admin 
    admin := chi.NewRouter()
    admin.Get("/metrics", apiCfg.HitHandler)
    // Mounting admin
    r.Mount("/admin", admin)
  

    // Serve Cors 
    srv := &http.Server {
        Addr: ":" + port,
        Handler: api.MiddlewareCors(r),
    }

    log.Printf("Serving files from %s on port localhost:%s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
    dbrequester.DeleteDB()
}

