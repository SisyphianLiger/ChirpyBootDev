package db_logic

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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
    
    

    userInfo := UserInfo{} 
    for key, Credentials := range dbToMem.Credentials {
        if Credentials.Email == email {
            if bcrypt.CompareHashAndPassword([]byte(Credentials.HashedPassword), []byte(password)) != nil {
                return UserInfo{}, errors.New("passwords do not match")
            }
            userInfo.Id = key
            userInfo.Email = Credentials.Email
            log.Printf("Login Credentials Successful for password %s", password)
            break
        }
    }
 
   
    return userInfo, nil
}

