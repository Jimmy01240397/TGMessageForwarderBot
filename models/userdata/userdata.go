package userdata

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/go-errors/errors"
    "MessageForwarder/utils/database"
)

type Userdata struct {
    ID string
    Uid string
    Chatid string
    Data map[string]any
}   

const tablename = "userdata"


func init() {
    _, err := database.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id varchar(255) NOT NULL PRIMARY KEY,
            uid varchar(255) NOT NULL,
            chatid varchar(255) NOT NULL,
            data varchar(255) NOT NULL DEFAULT "{}"
        )
    `, tablename))
    if err != nil {
        log.Panicln(err)
        return
    }
}

func Get(uid, chatid string) (*Userdata, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE uid=? and chatid=? LIMIT 1", tablename), uid, chatid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if !rows.Next() {
        return nil, nil
    }
    result := new(Userdata)
    var datajson string
    rows.Scan(&result.ID, &result.Uid, &result.Chatid, &datajson)
    json.Unmarshal([]byte(datajson), &result.Data)
    return result, nil
}

func Set(uid, chatid string, data map[string]any) (error) {
    datajson, err := json.Marshal(data)
    userdata, err := Get(uid, chatid)
    if err != nil {
        return err
    }
    if userdata == nil {
        _, err = database.Exec(fmt.Sprintf("INSERT INTO %s (id, uid, chatid, data) VALUES (?,?,?,?)", tablename), uuid.New().String(), uid, chatid, string(datajson))
    } else {
        _, err = database.Exec(fmt.Sprintf("UPDATE %s SET data=? where uid=? and chatid=?", tablename), string(datajson), uid, chatid)
    }
    return err
}

func Delete(uid, chatid string) (error) {
    userdata, err := Get(uid, chatid)
    if err != nil {
        return err
    }

    if userdata == nil {
        return errors.New("userdata not exist")
    }
    
    _, err = database.Exec(fmt.Sprintf("DELETE FROM %s where uid=? and chatid=?", tablename), uid, chatid)
    return err
}
