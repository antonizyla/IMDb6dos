package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"runtime"
)

var createTables = `
CREATE TABLE IF NOT EXISTS titles (
    tconst VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    titleType VARCHAR(32),
    primaryTitle VARCHAR(512),
    originalTitle VARCHAR(512),
    isAdult BOOLEAN,
    startYear INT,
    endYear INT,
    runtimeMinutes INT,
    genres TEXT[]
);

CREATE TABLE IF NOT EXISTS actors (
    nconst VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    primaryName VARCHAR(255),
    birthYear INT,
    deathYear INT,
    primaryProfession TEXT[],
    knownForTitles TEXT[]
);

CREATE TABLE IF NOT EXISTS title_actors (
    ID SERIAL PRIMARY KEY,
    tconst VARCHAR(255) NOT NULL,
    nconst VARCHAR(255) NOT NULL,
    title_ordering Integer,
    category TEXT, 
    job TEXT,
    characters TEXT,
    FOREIGN KEY (tconst) REFERENCES titles(tconst) ON DELETE CASCADE,
    FOREIGN KEY (nconst) REFERENCES actors(nconst) ON DELETE CASCADE
);
`

var db *sql.DB

func handleError(err error) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", filename, line, err)
	}
}

func init() {

	// initalise environment variables
	er := godotenv.Load()
	if er != nil {
		log.Fatalf("Error loading .env file")
	}

	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London", os.Getenv("HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), "5432")

	fmt.Println(connection)

	DB, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to database")
	}

	_, err = DB.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Created tables")
	}

	db = DB

}
