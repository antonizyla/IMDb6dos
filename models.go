package main

import "github.com/lib/pq"

type Title struct {
	Tconst         string         `db:"tconst" json:"tconst"`
	TitleType      string         `db:"titletype" json:"titletype"`
	PrimaryTitle   string         `db:"primarytitle" json:"primarytitle"`
	OriginalTitle  string         `db:"originaltitle" json:"originaltitle"`
	IsAdult        bool           `db:"isadult" json:"isadult"`
	StartYear      int            `db:"startyear" json:"startyear"`
	EndYear        int            `db:"endyear" json:"endyear"`
	RuntimeMinutes int            `db:"runtimeminutes" json:"runtimeminutes"`
	Genres         pq.StringArray `db:"genres" json:"genres"`
}

type Actor struct {
	Nconst            string         `db:"nconst" json:"nconst"`
	PrimaryName       string         `db:"primaryname" json:"primaryname"`
	BirthYear         int            `db:"birthyear" json:"birthyear"`
	DeathYear         int            `db:"deathyear" json:"deathyear"`
	PrimaryProfession pq.StringArray `db:"primaryprofession" json:"primaryprofession"`
	KnownForTitles    pq.StringArray `db:"knownfortitles" json:"knownfortitles"`
}

type MovieActor struct {
	ID             int    `db:"id" json:"id"`
	Tconst         string `db:"tconst" json:"tconst"`
	Title_ordering int    `db:"title_ordering" json:"title_ordering"`
	Nconst         string `db:"nconst" json:"nconst"`
	Category       string `db:"category" json:"category"`
	Job            string `db:"job" json:"job"`
	Characters     string `db:"characters" json:"characters"`
}
