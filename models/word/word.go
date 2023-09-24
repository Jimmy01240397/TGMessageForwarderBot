package word

import (
    "fmt"
    "log"
    "github.com/go-errors/errors"
    "github.com/google/uuid"
    "MessageForwarder/utils/database"
)

type Word struct {
    ID string
    Uid string
    Name string
    Word string
}   

const tablename = "word"


func init() {
    _, err := database.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id varchar(255) NOT NULL PRIMARY KEY,
            uid varchar(255) NOT NULL,
            name varchar(255) NOT NULL,
            word string NOT NULL DEFAULT ""
        )
    `, tablename))
    if err != nil {
        log.Panicln(err)
        return
    }
}

func Get(uid, name string) (*Word, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE uid=? and name=? LIMIT 1", tablename), uid, name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if !rows.Next() {
        return nil, nil
    }
    result := new(Word)
    rows.Scan(&result.ID, &result.Uid, &result.Name, &result.Word)
    return result, nil
}

func GetAll(uid string) ([]*Word, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE uid=?", tablename), uid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var words []*Word
    for rows.Next() {
        result := new(Word)
        rows.Scan(&result.ID, &result.Uid, &result.Name, &result.Word)
        words = append(words, result)
    }
    return words, nil
}

func Set(uid, name, word string) (error) {
    if (name == "") {
        return errors.New("name can't be empty")
    }
    worddata, err := Get(uid, name)
    if err != nil {
        return err
    }

    if worddata == nil {
        _, err = database.Exec(fmt.Sprintf("INSERT INTO %s (id, uid, name, word) VALUES (?,?,?,?)", tablename), uuid.New().String(), uid, name, word)
    } else {
        _, err = database.Exec(fmt.Sprintf("UPDATE %s SET word=? where uid=? and name=?", tablename), word, uid, name)
    }
    return err
}

func Delete(uid, name string) (error) {
    worddata, err := Get(uid, name)
    if err != nil {
        return err
    }

    if worddata == nil {
        return errors.New("word not exist")
    }
    
    _, err = database.Exec(fmt.Sprintf("DELETE FROM %s where uid=? and name=?", tablename), uid, name)
    return err
}
