package main

import "github.com/lib/pq"

type Title struct {
	Tconst         string   `db:"tconst"`
	TitleType      string   `db:"titletype"`
	PrimaryTitle   string   `db:"primarytitle"`
	OriginalTitle  string   `db:"originaltitle"`
	IsAdult        bool     `db:"isadult"`
	StartYear      int      `db:"startyear"`
	EndYear        int      `db:"endyear"`
	RuntimeMinutes int      `db:"runtimeminutes"`
	Genres         pq.StringArray `db:"genres"`
}

type Actor struct {
	Nconst            string         `db:"nconst"`
	PrimaryName       string         `db:"primaryname"`
	BirthYear         int            `db:"birthyear"`
	DeathYear         int            `db:"deathyear"`
	PrimaryProfession pq.StringArray `db:"primaryprofession"`
	KnownForTitles    pq.StringArray `db:"knownfortitles"`
}

type MovieActor struct {
	ID             int    `db:"id"`
	Tconst         string `db:"tconst"`
	Title_ordering int    `db:"title_ordering"`
	Nconst         string `db:"nconst"`
	Category       string `db:"category"`
	Job            string `db:"job"`
	Characters     string `db:"characters"`
}
