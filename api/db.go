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
    Users map[int]UserInfo `json:"email"`
    Credentials map[int]UserLogin `json:"password"`
}

// Generated Data from Chirp
type Chirp struct {
    Body string `json:"body"`
    ID int `json:"id"`
}

type UserLogin struct {
    HashedPassword string `json:"password"`
    Email string `json:"email"`
}

type UserInfo struct {
    Email string `json:"email"`
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
            Users: make(map[int]UserInfo), 
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

// Steps --> Make sure we check for duplicat emails
// Make sure we pass 
func (db *DBrequester) CreateUser(email string, password string) (UserInfo, error) {

    err := db.ensureDB()
   
    if err != nil {
        log.Print(err)
        return UserInfo{}, err
    }


    // hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
   
    if err != nil {
        log.Print(err)
        return UserInfo{}, err
    }
   

    userLogin := UserLogin {
            Email: email,
            HashedPassword: password,
    }

    userInfo := UserInfo{
        Email: email,
        ID: 0,
    }

    if err != nil {
        _, err := NewDB(db.path)     
        if err != nil {
            log.Printf("Something went wrong creating the DB check path string")
            return UserInfo{}, nil
        }

        dbToMem, _ := db.loadDB()

        err = db.NoRepeatEmails(email, &dbToMem)

        if err != nil {
            return UserInfo{}, err
        }

        dbToMem.Users[0] = userInfo
        dbToMem.Credentials[0] = userLogin
        err = db.writeDB(dbToMem)
        if err != nil {
            log.Printf("Failed to Write to DB, path may be corrupt")
            return UserInfo{}, err
        }
        return UserInfo{}, err
    }

    dbToMem, _ := db.loadDB()
    err = db.NoRepeatEmails(email, &dbToMem)

    if err != nil {
        return UserInfo{}, err
    }

    nextAdd := len(dbToMem.Users) + 1
    userInfo.ID = nextAdd
    dbToMem.Users[nextAdd] = userInfo
    dbToMem.Credentials[nextAdd] = userLogin
    resp := dbToMem.Users[nextAdd] 
    err = db.writeDB(dbToMem)
    if err != nil {
        log.Printf("Cannot put struct into db check Chirp body")
        return UserInfo{}, err
    }
    return resp, nil
}

func (db *DBrequester) NoRepeatEmails(email string, dbToMem *DBStructure) error {
 
    for _, user := range dbToMem.Users {
        if email == user.Email {
            return errors.New("Email already in use, please sign up with another email")
        }
    }
    return nil
}



func (db *DBrequester) Login(email string, password string) (UserInfo, error){
     
    err := db.ensureDB()
   
   
    if err != nil {
        log.Print(err)
        return UserInfo{}, err
    }

    dbToMem, err := db.loadDB()

    if err != nil {
        log.Printf("Failed to Write to DB, path may be corrupt")
        return UserInfo{}, err
    }
    



    var user UserInfo 
    for _, Credentials := range dbToMem.Credentials {
        log.Printf("Made it to Credentials Check")
        if Credentials.Email == email {
            user, err = findEmail(email, &dbToMem) 
            if err != nil {
                return UserInfo{}, err
            }
            if Credentials.HashedPassword != password {
                return UserInfo{}, errors.New("passwords do not match")
            }
            break
        }
    }

    // err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
    
    
    if err != nil {
        log.Print("password does not match login")
        return UserInfo{}, nil
    }

   
    return user, nil
}


func findEmail(email string, dbToMem *DBStructure) (UserInfo, error) {

    for _, userEmail := range dbToMem.Users {
        if userEmail.Email == email {
            return userEmail, nil
        }
    }

    return UserInfo{}, errors.New("Email Not Found")
}
