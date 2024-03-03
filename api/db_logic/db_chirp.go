package db_logic


import (
	"fmt"
	"log"
	"sort"
    "strings"
)

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
