package database

import (
    "log"
    "os"
    "sync"
//    "fmt"
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
    "github.com/go-errors/errors"

    "MessageForwarder/utils/config"
)

var db *sql.DB
var lock *sync.RWMutex

func init() {
    lock = new(sync.RWMutex)
    lock.Lock()
    defer lock.Unlock()
    var err error
    if _, err := os.Stat(config.DBname); errors.Is(err, os.ErrNotExist) {
        os.Create(config.DBname)
    }
    db, err = sql.Open("sqlite3", config.DBname)
    if err != nil {
        log.Panicln(err)
    }
}

func Exec(sqlcmd string, args ...any) (sql.Result, error) {
    lock.Lock()
    defer lock.Unlock()
    stmt, err := db.Prepare(sqlcmd)
    if err != nil {
        return nil, err
    }
    return stmt.Exec(args...)
}

func Query(sqlcmd string, args ...any) (*sql.Rows, error) {
    lock.RLock()
    defer lock.RUnlock()
    stmt, err := db.Prepare(sqlcmd)
    if err != nil {
        return nil, err
    }
    return stmt.Query(args...)
}

func Close() {
    lock.Lock()
    defer lock.Unlock()
    db.Close()
}
