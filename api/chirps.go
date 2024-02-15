package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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



func Valid_Post_Chirp(w http.ResponseWriter, r *http.Request){
    
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
    // Gotta fix dis
    cResp := Chirp{Body: body, id: 0}
        respJaysawn, _ := json.Marshal(cResp)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(201)
        w.Write(respJaysawn)

}







