package main 

import "gorm.io/gorm"

type book struct {
	ID     int `gorm:"primaryKey"`
	Title  string
	Author string
	Isb    bool
	End    int
}

type Title struct {
    gorm.Model
    Name string
    Genres         []Genre `gorm:"many2many:title_genres;"`
}

type Genre struct {
    gorm.Model 
    Genre string 
}
