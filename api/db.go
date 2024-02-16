package api

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
    "strings" 
    "net/http"
)

// Used to lock DB when accessing
type DBrequester struct {
    path string
    mux *sync.RWMutex
}

// Access of DB
type DBStructure struct {
    Chirps map[int]Chirp `json:"chirps"`
}

// Generated Data from Chirp
type Chirp struct {
    Body string `json:"cleaned_body"`
    id int `json:"id"`
}


func Filter_Post(body string) string {

    filter := strings.Split(body, " ")
    for i,word := range filter {
        tmp := strings.ToLower(word)
        if tmp == "kerfuffle" || tmp == "sharbert" || tmp == "fornax" {
            filter[i] = "****"            
        }
    }
    return strings.Join(filter, " ")
}



func (db *DBrequester) CreateChirp(w http.ResponseWriter, r *http.Request){
   

    // Decoded Body inc
    type Request struct {
        Body string `json:"body"`
    }

    // We want to process the incoming body 
    decoder := json.NewDecoder(r.Body)

    requests := Request{}
    err := decoder.Decode(&requests)
   
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(500)
        return
    }


    // From this point 
    // Want to put response to body
    type ErrorResp struct {
        ErrorResp string `json:"error"`
    }

    // Check Length
    if len(requests.Body) > 140 {
        err := ErrorResp{ErrorResp: "Chirp is too long"}
        errJson, _ := json.Marshal(err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400) 
        w.Write(errJson)
        return
    }


    // Cursed Case 
    body := Filter_Post(requests.Body)

    // so load db now I guess?

    dbToMem := DBStructure{}
    if db.ensureDB() != nil {
        db.mux.Lock()
        NewDB(db.path)
        // Insert a new dbToMem here 
        newChirp := Chirp{Body:body, id:0} 
        dbToMem.Chirps[newChirp.id] = newChirp
        db.writeDB(dbToMem) 
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(201)
        w.Write([]byte(body))
        db.mux.Unlock()
    } else {
        // Count 
        db.mux.Lock()
        dbToMem ,err = db.loadDB()
        cnt := len(dbToMem.Chirps)
        input := DBStructure{}
        newChirp := Chirp{Body:body, id:cnt}
        input.Chirps[cnt+1] = newChirp
        db.writeDB(input)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(201)
        w.Write([]byte(body))
        db.mux.Unlock()
    }
}


func (db *DBrequester) GetChirp(w http.ResponseWriter, r *http.Request){
    
    // Decoded Body inc
    type Request struct {
        Body string `json:"body"`
    }

    // We want to process the incoming body 
    decoder := json.NewDecoder(r.Body)

    requests := Request{}
    err := decoder.Decode(&requests)
   
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(500)
        return
    }


    // From this point 
    // Want to put response to body
    type ErrorResp struct {
        ErrorResp string `json:"error"`
    }

    // Check Length
    if len(requests.Body) > 140 {
        err := ErrorResp{ErrorResp: "Chirp is too long"}
        errJson, _ := json.Marshal(err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400) 
        w.Write(errJson)
        return
    }


    // Cursed Case 
    body := Filter_Post(requests.Body)
    cResp := Chirp{Body: body, id: 0}
        respJaysawn, _ := json.Marshal(cResp)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(201)
        w.Write(respJaysawn)

}

// Generates a new DB is no DB is made
func NewDB(path string) (*DBrequester, error) {

    // Check if file does not exist
    _, err :=os.ReadFile("database.json")

    // Now there is an error which means we make the db and start the struct
    if err != nil { 
        os.WriteFile(path, []byte{}, 0600)
        req := DBrequester{path:path, mux:&sync.RWMutex{}}
        return &req, err
    }
    
    req := DBrequester{path:path, mux:&sync.RWMutex{}}
    return &req, err
}


func (db *DBrequester) loadDB() (DBStructure, error) {
    // Read the JSON file into a byte array  
    jsonData, err := os.ReadFile(db.path)

    // If we request when there is no DB we return error
    if err != nil {
        log.Println(err) 
        return DBStructure{}, err
    }
    
    dbToMem := DBStructure{}

    // Unmarshal JSON data
    err = json.Unmarshal(jsonData, &dbToMem)

    // If unmarshal goes wrong to struct
    if err != nil {
        log.Println(err)
        return DBStructure{}, err
    }
    // we have struct 
    return dbToMem, nil
}


func (db *DBrequester) ensureDB() error { 
    _, err := os.ReadFile(db.path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist){
            err := os.WriteFile(db.path, []byte{}, 0600)
            return err 
        }
        return err
    }
    return nil
}


// writeDB writes the database file to disk
func (db *DBrequester) writeDB(dbStructure DBStructure) error {
    sendingBody, err := json.Marshal(&dbStructure)

    os.WriteFile(db.path, sendingBody, 0600) 
    if err != nil {
        log.Println(err)
        return err 
    }
        return nil
}





