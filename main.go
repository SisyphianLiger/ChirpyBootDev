package main

import (
	"log"
	"net/http"
	// Using dot convention to becasue of clear Handler names
	"github.com/SisyphianLiger/Chirpy/api"
)


func main() {

    const filepathRoot = "./app"
    const port = "8000"

    fs := http.FileServer(http.Dir(filepathRoot))
    // Use the http.NewServeMux() function to create an empty servemux
    mux := http.NewServeMux()
    
    // Creating a empty struct to track state of Handler
    apiCfg := &api.ApiConfig{}
    mux.Handle("/app", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", fs)))
    mux.HandleFunc("/healthz", api.HealthzHandler)
    mux.HandleFunc("/metrics", apiCfg.HitHandler)
    mux.HandleFunc("/reset", apiCfg.ResetHandler)

    // Serve Cors 
    srv := &http.Server {
        Addr: ":" + port,
        Handler: api.MiddlewareCors(mux),
    }

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

