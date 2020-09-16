package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	initDB()
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/checktoken", CheckToken)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func initDB() {
	var err error
	os.Create("./foo.db")
	db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS userinfo (username TEXT,password TEXT)")
	if err != nil {
		panic(err)
	}
	stmt.Exec()

}
