package main

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
	Genres         []Genre `gorm:"many2many:title_genres;"`
}

type Genre struct {
	ID     int `gorm:"primaryKey;autoIncrement;unique"`
	Genre  string `gorm:"unique"`
}
