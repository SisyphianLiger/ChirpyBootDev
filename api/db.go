package api

import (
	"encoding/json"
	"log"
	"os"
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


func (db *DBrequester) LoadDB() (DBStructure, error) {
    // Read the JSON file into a byte array  
    jsonData, err := os.ReadFile(db.path)

    // If we request when there is no DB we return error
    if err != nil {
        log.Fatal(err) 
        return DBStructure{}, err
    }
    
    dbToMem := DBStructure{}

    // Unmarshal JSON data
    err = json.Unmarshal(jsonData, &dbToMem)

    if err != nil {
        return DBStructure{}, err
    }
    // we have struct 
    return dbToMem, err




}
