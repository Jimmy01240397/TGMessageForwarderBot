package rule

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/go-errors/errors"
    "github.com/google/uuid"
    "MessageForwarder/utils/database"
)

type Rule struct {
    ID string
    Uid string
    Name string
    ChatIDlist []string
    Whitelist string
    Blacklist string
    Replacelist []Replace
}   

type Replace struct {
    From string `json:"from"`
    To string `json:"to"`
}

const tablename = "rule"


func init() {
    _, err := database.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id varchar(255) NOT NULL PRIMARY KEY,
            uid varchar(255) NOT NULL,
            name varchar(255) NOT NULL,
            chatidlist string NOT NULL DEFAULT "[]",
            whitelist string NOT NULL DEFAULT "",
            blacklist string NOT NULL DEFAULT "",
            replacelist string NOT NULL DEFAULT "[]"
        )
    `, tablename))
    if err != nil {
        log.Panicln(err)
        return
    }
}

func Get(uid, name string) (*Rule, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE uid=? and name=? LIMIT 1", tablename), uid, name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if !rows.Next() {
        return nil, nil
    }
    result := new(Rule)
    var chatidlistjson string
    var replacelistjson string
    rows.Scan(&result.ID, &result.Uid, &result.Name, &chatidlistjson, &result.Whitelist, &result.Blacklist, &replacelistjson)
    json.Unmarshal([]byte(chatidlistjson), &result.ChatIDlist)
    json.Unmarshal([]byte(replacelistjson), &result.Replacelist)
    return result, nil
}

func GetAll(uid string) ([]*Rule, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE uid=?", tablename), uid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var rules []*Rule
    for rows.Next() {
        result := new(Rule)
        var chatidlistjson string
        var replacelistjson string
        rows.Scan(&result.ID, &result.Uid, &result.Name, &chatidlistjson, &result.Whitelist, &result.Blacklist, &replacelistjson)
        json.Unmarshal([]byte(chatidlistjson), &result.ChatIDlist)
        json.Unmarshal([]byte(replacelistjson), &result.Replacelist)
        rules = append(rules, result)
    }
    return rules, nil
}

func Set(uid, name string, chatidlist []string, whitelist, blacklist string, replacelist []Replace) (error) {
    if (name == "") {
        return errors.New("name can't be empty")
    }
    chatidlistjson, err := json.Marshal(chatidlist)
    replacelistjson, err := json.Marshal(replacelist)
    rule, err := Get(uid, name)
    if err != nil {
        return err
    }

    if rule == nil {
        _, err = database.Exec(fmt.Sprintf("INSERT INTO %s (id, uid, name, chatidlist, whitelist, blacklist, replacelist) VALUES (?,?,?,?,?,?,?)", tablename), uuid.New().String(), uid, name, string(chatidlistjson), whitelist, blacklist, string(replacelistjson))
    } else {
        _, err = database.Exec(fmt.Sprintf("UPDATE %s SET chatidlist=?, whitelist=?, blacklist=?, replacelist=? where uid=? and name=?", tablename), string(chatidlistjson), whitelist, blacklist, string(replacelistjson), uid, name)
    }
    return err
}

func Delete(uid, name string) (error) {
    rule, err := Get(uid, name)
    if err != nil {
        return err
    }

    if rule == nil {
        return errors.New("rule not exist")
    }
    
    _, err = database.Exec(fmt.Sprintf("DELETE FROM %s where uid=? and name=?", tablename), uid, name)
    return err
}
