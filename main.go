package main

import (
    "net/http"
)

type apiHandler struct{}


func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
   
    // Use the http.NewServeMux() function to create an empty servemux
    mux := http.NewServeMux()
    mux.Handle("/", http.FileServer(http.Dir(".")))
    corsMux := middlewareCors(mux)
    srv := &http.Server {
        Addr: "localhost:3000",
        Handler: corsMux,
    }
    err := srv.ListenAndServe();
    if err != nil {
        panic(err);
    }

}
