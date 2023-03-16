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
    // if seed is selected get all data
	if *seedptr {
		fmt.Println("Getting Data from IMDB")
		DispatchDownloads()
		DispatchInsertions()
	}

    findPath("nm0000206", "nm0000124")
    

}
