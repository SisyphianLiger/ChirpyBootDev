package api

import (
	"encoding/json"
	"net/http"
)


func HealthzHandler( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (db *DBrequester) GetChirps(w http.ResponseWriter, r *http.Request) {
    dbresp, err := db.GetChirp()

    if err != nil {
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
    } else {
        w.WriteHeader(200)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(dbresp)
    }
}

func (db *DBrequester) PostChirps( w http.ResponseWriter, r *http.Request) {
          
    var req Request
    dec := json.NewDecoder(r.Body)
    err := dec.Decode(&req)
   
    if err != nil {
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
        return
    }        

    newChirp, err := db.CreateChirp(req.Body)

    if err != nil {
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
        return
    }        


    w.WriteHeader(201)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newChirp)

    


}
