package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"    
    "github.com/joho/godotenv"
)

type Title struct {
	gorm.Model
	tconst         string
	titleType      string
	primaryTitle   string
	originalTitle  string
	isAdult        bool
	startYear      int
	endYear        int
	runtimeMinutes int
	genres         []string
}

func test() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    db, err := gorm.Open(postgres.Open("host=localhost user="+os.Getenv("POSTGRES_USER")+" password="+os.Getenv("POSTGRES_PASSWORD")+" dbname="+os.Getenv("POSTGRES_DB")+ " port=5432 sslmode=disable "), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Title{})
}

func read_file(DItem DonwloadItem) {

}
