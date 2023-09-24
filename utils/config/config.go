package config

import (
    "os"
    "log"
    "strings"
    "github.com/joho/godotenv"
)

var Token string
var DBname string
var Debug bool

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Panicln("Error loading .env file")
    }
    Token = os.Getenv("TGBOTTOKEN")
    exists := false
    DBname, exists = os.LookupEnv("DBNAME")
    if !exists {
        DBname = "data.db"
    }
    Debug = strings.ToLower(os.Getenv("DEBUG")) == "true"
}

