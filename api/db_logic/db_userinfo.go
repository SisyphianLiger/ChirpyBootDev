package db_logic

import (
    "errors"
    "log"
	"golang.org/x/crypto/bcrypt"

)


// Steps --> Make sure we check for duplicat emails
// Make sure we pass 
func (db *DBrequester) CreateUser(email string, password string) (UserInfo, error) {

    err := db.ensureDB()
   
    if err != nil {
        log.Print(err)
        return UserInfo{}, err
    }


    hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
   
    if err != nil {
        log.Print(err)
        return UserInfo{}, err
    }
   

    userLogin := UserLogin {
            Email: email,
            HashedPassword: hashedPw,
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

    nextAdd := len(dbToMem.Credentials) + 1
    dbToMem.Credentials[nextAdd] = userLogin
    resp := UserInfo{Id:nextAdd, Email: email}

    err = db.writeDB(dbToMem)
    if err != nil {
        log.Printf("Cannot put struct into db check Chirp body")
        return UserInfo{}, err
    }
    return resp, nil
}

func (db *DBrequester) NoRepeatEmails(email string, dbToMem *DBStructure) error {
 
    for _, user := range dbToMem.Credentials {
        if email == user.Email {
            return errors.New("Email already in use, please sign up with another email")
        }
    }
    return nil
}
