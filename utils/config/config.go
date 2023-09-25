package config

import (
    "os"
    //"log"
    "strings"
    "path/filepath"
    "github.com/joho/godotenv"
)

var Token string
var DBname string
var Debug bool

func init() {
    ex, err := os.Executable()
    if err == nil {
        exPath := filepath.Dir(ex)
        os.Chdir(exPath)
    }
    err = godotenv.Load()
    Token = os.Getenv("TGBOTTOKEN")
    exists := false
    DBname, exists = os.LookupEnv("DBNAME")
    if !exists {
        DBname = "data.db"
    }
    Debug = strings.ToLower(os.Getenv("DEBUG")) == "true"
}

