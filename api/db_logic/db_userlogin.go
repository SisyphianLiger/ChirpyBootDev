package db_logic

import (
    "log"
    "errors"
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
            if Credentials.HashedPassword != password {
                return UserInfo{}, errors.New("passwords do not match")
            }
            userInfo.Id = key
            userInfo.Email = Credentials.Email
            break
        }
    }

 
    
    if err != nil {
        log.Print("password does not match login")
        return UserInfo{}, nil
    }
    
   
    return userInfo, nil
}

