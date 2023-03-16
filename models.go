package main

type Title struct {
	ID             int
	Tconst         string
	TitleType      string
	PrimaryTitle   string
	OriginalTitle  string
	IsAdult        bool
	StartYear      int
	EndYear        int
	RuntimeMinutes int
	Genres         []string
}

type Actor struct {
	ID                int
	Nconst            string
	PrimaryName       string
	BirthYear         int
	DeathYear         int
	PrimaryProfession []string
	KnownForTitles    []string
}

type MovieActor struct {
	ID         int
	Tconst     string
	Ordering   int
	Nconst     string
	Category   string
	Job        string
	Characters string
}
