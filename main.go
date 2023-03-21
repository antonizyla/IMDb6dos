package main

/*
 #cgo LDFLAGS: -L. / -lcfile
 #include "code.h"
*/

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Env struct {
	db *sql.DB
}

func main() {

	// Parse command line flags
	seedptr := flag.Bool("seed", false, "Seed the database with some data")
	clearDBptr := flag.Bool("clear", false, "Reset Database")

	flag.Parse()
	// if seed is selected get all data
	if *seedptr {
		fmt.Println("Getting Data from IMDB")
		DispatchDownloads()
		DispatchInsertions()
	} else if *clearDBptr {
		fmt.Println("Clearing Database")
		_, err := db.Exec("Drop database $1", os.Getenv("POSTGRES_DB"))
		handleError(err)
		fmt.Println("Database Cleared")
		fmt.Println("Run Program with -seed to seed the database")
	}

	http.HandleFunc("/api/actors/info", actorInfoHandler)
	http.HandleFunc("/api/titles/info", titlesInfoHandler)
	http.ListenAndServe(":8080", nil)
}
