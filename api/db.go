package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
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
    Body string `json:"body"`
    ID int `json:"id"`
}

// Decoded Body inc
type Request struct {
    Body string `json:"body"`
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



func (db *DBrequester) CreateChirp(body string) (Chirp, error){
  
    // Check Length
    if len(body) > 140 {
        log.Printf("Error: body is to long to post")
        return Chirp{}, fmt.Errorf("Error: body is two long to post")
    }
    // Cursed Case 
    postBody := Filter_Post(body)

    err := db.ensureDB()
   
    newChirp := Chirp{
        ID: 0,
        Body: postBody,
    }
    
    if err != nil {
        _, err := NewDB(db.path)     
        if err != nil {
            log.Printf("Something went wrong creating the DB check path string")
            return Chirp{}, nil
        }
        dbToMem, _ := db.loadDB()
        dbToMem.Chirps[0] = newChirp
        err = db.writeDB(dbToMem)
        if err != nil {
            log.Printf("Failed to Write to DB, path may be corrupt")
            return Chirp{}, err
        }
        return Chirp{}, err
    }

    dbToMem, _ := db.loadDB()
    nextAdd := len(dbToMem.Chirps) + 1
    newChirp.ID = nextAdd
    dbToMem.Chirps[nextAdd] = newChirp
    resp := dbToMem.Chirps[nextAdd] 
    err = db.writeDB(dbToMem)
    if err != nil {
        log.Printf("Cannot put struct into db check Chirp body")
        return Chirp{}, err
    }
    return resp, nil
}


func (db *DBrequester) GetChirp() ([]Chirp, error){
  
    err := db.ensureDB()

    if err != nil {
        log.Printf("No DB Found")
        return []Chirp{}, err
    }

    dbToMem, err :=  db.loadDB()
    if err != nil {
        log.Printf("DB not loaded successfully")
        return []Chirp{}, err
    }
   
   
    // Need an array here sorted, 
    // Need to return an array 
    var dbresp []Chirp

    fmt.Print(dbToMem.Chirps)

    for _, chirp := range dbToMem.Chirps {
        dbresp = append(dbresp, chirp)
    }

    sort.Slice(dbresp, func(i, j int) bool {
        return dbresp[i].ID < dbresp[j].ID
    })

    return dbresp, nil

}


func (db *DBrequester) GetChirpID(id int) (Chirp, error){
  
    err := db.ensureDB()

    if err != nil {
        log.Printf("No DB Found")
        return Chirp{}, err
    }

    dbToMem, err :=  db.loadDB()
    if err != nil {
        log.Printf("DB not loaded successfully")
        return Chirp{}, err
    }
   
   
    // Need an array here sorted, 
    // Need to return an array 
    var chirpFromid Chirp

    for _, chirp := range dbToMem.Chirps {
        fmt.Printf("Current Chirp id %v\n", chirp.ID)
        if chirp.ID == id {
            chirpFromid = chirp
            break
        }
    }

    if chirpFromid.ID != id {
        return Chirp{}, fmt.Errorf("Chirp with ID %v not found", id)
    }

    return chirpFromid, nil

}



// Generates a new DB is no DB is made
func NewDB(path string) (*DBrequester, error) {

    // Check if file does not exist
    _, err :=os.ReadFile("database.json")

    // Now there is an error which means we make the db and start the struct
    if err != nil { 
        initialDB := DBStructure{
            Chirps: make(map[int]Chirp), 
        }

        jsonData, err := json.Marshal(initialDB)

        if err != nil {
            log.Println("Failed to marshal initial DB Structure")
            return nil, err 
        }

        err = os.WriteFile(path, jsonData, 0600)
        if err != nil {
            log.Printf("Could not create database, please check path")
        }
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
    if err != nil {
        log.Println(err)
        return err 
    }

    err = os.WriteFile(db.path, sendingBody, 0600) 
    if err != nil {
        log.Println(err)
        return err 
    }
        return nil
}

func (db *DBrequester) DeleteDB() error{
    
    err := db.ensureDB()
    if err != nil {
        return err
    }
    
    err = os.Remove(db.path)
    if err != nil {
        log.Printf("Error removing db check path")
        return err
    }
    
    var newDB DBStructure
    err = db.writeDB(newDB)

    if err != nil {
        log.Printf("Could not create the db try running writeDB again")
        return err
    }
    
    return nil

}



