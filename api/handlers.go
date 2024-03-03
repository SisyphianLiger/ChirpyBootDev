package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
    "github.com/SisyphianLiger/Chirpy/api/db_logic"
)

// Used to access db from db_logic
type DBreq struct {
    *db_logic.DBrequester
}

type LoginRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

func HealthzHandler( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, err := w.Write([]byte(http.StatusText(http.StatusOK)))
    if err != nil {
        log.Print(err)
        return 
    }
}


func (db *DBreq) HandleGetID(w http.ResponseWriter, r *http.Request) {
    

    idStr := chi.URLParam(r, "chirpID")
    ID, err := strconv.Atoi(idStr)

    if err != nil {
        log.Printf("Invalid ID in URL")
        w.WriteHeader(404)
        w.Header().Set("Content-Type", "application/json")
        return 
    }

    chirp, err := db.GetChirpID(ID)
   
    if err != nil {
        log.Printf("ID Does not exist")
        w.WriteHeader(404)
        w.Header().Set("Content-Type", "application/json")
        return
    }        


    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")

    err = json.NewEncoder(w).Encode(chirp) 
    if err != nil {
        return
    }
}


func (db *DBreq) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
    dbresp, err := db.GetChirp()

    if err != nil {
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
    } else {
        w.WriteHeader(200)
        w.Header().Set("Content-Type", "application/json")
        err = json.NewEncoder(w).Encode(dbresp)
        if err != nil {
            return
        }
    }
}

func (db *DBreq) HandlePostChirp( w http.ResponseWriter, r *http.Request) {
    
    var req db_logic.Request
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
    err = json.NewEncoder(w).Encode(newChirp)
    if err != nil {
        return
    }

}


func (db *DBreq) HandlePostUser( w http.ResponseWriter, r *http.Request) {

    var req LoginRequest 
    dec := json.NewDecoder(r.Body)

    err := dec.Decode(&req)
   
    if err != nil {
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
        return
    }        

    // Here is where we need a nother way to handle 
    createUser, err := db.CreateUser(req.Email, req.Password)

    if err != nil {
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
        return
    }        


    w.WriteHeader(201)
    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(createUser)
    if err != nil {
        return
    }
    
}

func (db *DBreq) HandlePostLogin( w http.ResponseWriter, r *http.Request) {
    var req LoginRequest 
    dec := json.NewDecoder(r.Body)
    err := dec.Decode(&req)

    if err != nil {
        w.WriteHeader(401)
        w.Header().Set("Content-Type", "application/json")
        return
    }        

    // Here is where we need a nother way to handle 
    login, err := db.Login(req.Email, req.Password)

    if err != nil {
        w.WriteHeader(401)
        w.Header().Set("Content-Type", "application/json")
        return
    }        


    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(login)
    if err != nil {
        return
    }
}
