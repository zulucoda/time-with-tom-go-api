package database

import (
    "database/sql"
    "log"
)

var DbConn *sql.DB

func SetupDatabase() {
    var err error
    DbConn, err = sql.Open("mysql", "root:sa@tcp(127.0.0.1:3306)/time_with_tom_db")
    if err != nil {
        log.Fatal(err)
    }
}