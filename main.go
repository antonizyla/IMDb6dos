package main

import (
	"database/sql"
	"flag"
	"fmt"
)

type Env struct {
	db *sql.DB
}

func main() {

	// Parse command line flags
	seedptr := flag.Bool("seed", false, "Seed the database with some data")
	flag.Parse()

	if *seedptr {
		fmt.Println("Getting Data from IMDB")
		DispatchDownloads()
		DispatchInsertions()
	}

}
