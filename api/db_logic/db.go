package db_logic

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
    // "golang.org/x/crypto/bcrypt"
)

// Used to lock DB when accessing
type DBrequester struct {
    path string
    mux *sync.RWMutex
}

// Access of DB
type DBStructure struct {
    Chirps map[int]Chirp `json:"chirps"`
    Credentials map[int]UserLogin `json:"password"`
}

// Generated Data from Chirp
type Chirp struct {
    Body string `json:"body"`
    ID int `json:"id"`
}

type UserInfo = struct {
    Email string `json:"email"`
    Id int `json:"id"`
}

type UserLogin struct {
    HashedPassword []byte `json:"password"`
    Email string `json:"email"`
}

// Decoded Body inc
type Request struct {
    Body string `json:"body"`
}




// Generates a new DB is no DB is made
func NewDB(path string) (*DBrequester, error) {

    // Check if file does not exist
    _, err :=os.ReadFile("database.json")

    // Now there is an error which means we make the db and start the struct
    if err != nil { 
        initialDB := DBStructure{
            Chirps: make(map[int]Chirp), 
            Credentials: make(map[int]UserLogin),
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




