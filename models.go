package main

import "github.com/lib/pq"

type book struct {
	ID     int `gorm:"primaryKey"`
	Title  string
	Author string
	Isb    bool
	End    int
}

type Title struct {
	ID             int `gorm:"primaryKey"`
	Tconst         string
	TitleType      string
	PrimaryTitle   string
	OriginalTitle  string
	IsAdult        bool
	StartYear      int
	EndYear        int
	RuntimeMinutes int
	Genres         pq.StringArray `gorm:"type:text[]"`
}

type Actor struct {
    ID               int `gorm:"primaryKey"`
    Nconst           string
    PrimaryName      string
    BirthYear        int
    DeathYear        int
    PrimaryProfession pq.StringArray `gorm:"type:text[]"`
    KnownForTitles   pq.StringArray `gorm:"type:text[]"`
}
