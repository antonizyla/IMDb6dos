package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

var createTables = `
CREATE TABLE IF NOT EXISTS titles (
    tconst VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    titleType VARCHAR(255),
    primaryTitle VARCHAR(255),
    originalTitle VARCHAR(255),
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

CREATE TABLE IF NOT EXISTS movie_actors (
    id SERIAL PRIMARY KEY,
    tconst VARCHAR(255) NOT NULL,
    nconst VARCHAR(255) NOT NULL,
    ordering INT,
    category VARCHAR(255),
    job VARCHAR(255),
    characters TEXT[]
    FOREIGN KEY (tconst) REFERENCES titles(tconst) ON DELETE CASCADE,
    FOREIGN KEY (nconst) REFERENCES actors(nconst) ON DELETE CASCADE
);
`


func init() {

	// initalise environment variables
	er := godotenv.Load()
	if er != nil {
		log.Fatalf("Error loading .env file")
	}

	// create a connection pool to the database
	var err error
	db, err = gorm.Open(postgres.Open("host=localhost user="+os.Getenv("POSTGRES_USER")+" password="+os.Getenv("POSTGRES_PASSWORD")+" dbname="+os.Getenv("POSTGRES_DB")+" port=5432 sslmode=disable "), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)
	}

	// initalise raw db connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// create pool of connections
	sqlDB.SetConnMaxIdleTime(30 * time.Second)
	sqlDB.SetMaxIdleConns(50)

	// create tables in the database from models
	db.AutoMigrate(&Title{}, &Actor{}, &MovieActor{})

	// db.Exec("ALTER TABLE titles ADD UNIQUE (tconst);")
	// db.Exec("ALTER TABLE actors ADD UNIQUE (nconst);")

	db.Exec("ALTER TABLE movie_actors ADD FOREIGN KEY (tconst) REFERENCES titles(tconst) ON DELETE CASCADE;")

	db.Exec("ALTER TABLE movie_actors ADD FOREIGN KEY (nconst) REFERENCES actors(nconst) ON DELETE CASCADE;")
}
