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
	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetMaxIdleConns(20)

	// create tables in the database from models
	db.AutoMigrate(&Title{}, &Genre{})

}
