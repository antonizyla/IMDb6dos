package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/lib/pq"
)

type DonwloadItem struct {
	url      string
	filename string
}

func DispatchDownloads() {

	os.Mkdir("data", 0777)
	var urls = [5]DonwloadItem{{"https://datasets.imdbws.com/title.ratings.tsv.gz", "data/title.ratings.tsv.gz"},
		{"https://datasets.imdbws.com/title.basics.tsv.gz", "data/title.basics.tsv.gz"}, {"https://datasets.imdbws.com/title.principals.tsv.gz", "data/title.principals.tsv.gz"}, {"https://datasets.imdbws.com/title.akas.tsv.gz", "data/title.akas.tsv.gz"}, {"https://datasets.imdbws.com/name.basics.tsv.gz", "data/name.basics.tsv.gz"}}

	var workGroup sync.WaitGroup
	for _, url := range urls {
		// check for if the file exists already
		_, error := os.Stat(url.filename[:len(url.filename)-3])
		if errors.Is(error, os.ErrNotExist) {
			workGroup.Add(1)
			go ProcessItem(url, &workGroup)
		} else {
			fmt.Println("File already exists: ", url.filename[:len(url.filename)-3])
		}
	}
	workGroup.Wait()
}

func DispatchInsertions() {
	// titles and actors can be done in parallel
	wg := sync.WaitGroup{}
	wg.Add(2)
	go InsertDatabaseTitlesWrap(&wg)
	go InsertDatabaseActorsWrap(&wg)
	wg.Wait()
	// insert link between titles and actors
	insertDatabaseMoviesActors()

	// delete all data files to free up space
	os.RemoveAll("data")
}

func InsertDatabaseTitlesWrap(wg *sync.WaitGroup) {
	defer wg.Done()
	InsertDatabaseTitles()
}

func InsertDatabaseActorsWrap(wg *sync.WaitGroup) {
	defer wg.Done()
	InsertDatabaseActors()
}

func read(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

func InsertDatabaseTitles() {
	path := "data/title.basics.tsv"

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	handleError(err)
	defer file.Close()

	fmt.Println("Inserting titles into database")

	transaction, err := db.Begin()
	handleError(err)

	stmt, err := transaction.Prepare(pq.CopyIn("titles", "tconst", "titletype", "primarytitle", "originaltitle", "isadult", "startyear", "endyear", "runtimeminutes", "genres"))
	handleError(err)

	reader := bufio.NewReader(file)
	i := 0
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if i > 0 {
			title := parseTitle(string(line))
			_, err = stmt.Exec(title.Tconst, title.TitleType, title.PrimaryTitle, title.OriginalTitle, title.IsAdult, title.StartYear, title.EndYear, title.RuntimeMinutes, pq.Array(title.Genres))
			handleError(err)
		}
		i++
	}
	_, err = stmt.Exec()
	handleError(err)

	err = stmt.Close()
	handleError(err)

	err = transaction.Commit()
	handleError(err)

	fmt.Println("Inserted titles")

}

func parseTitle(line string) Title {
	arr := strings.Split(line, "\t")
	tconst := arr[0]
	titleType := arr[1]
	primaryTitle := arr[2]
	originalTitle := arr[3]
	isAdult, _ := strconv.ParseBool(arr[4])
	startYear, _ := strconv.Atoi(arr[5])
	endYear, _ := strconv.Atoi(arr[6])
	runtimeMinutes, _ := strconv.Atoi(arr[7])
	genres := strings.Split(arr[8], ",")

	if len(genres) == 1 && genres[0] == "\\N" {
		genres = []string{"No Genre"}
	}

	return Title{
		Tconst:         tconst,
		TitleType:      titleType,
		PrimaryTitle:   primaryTitle,
		OriginalTitle:  originalTitle,
		IsAdult:        isAdult,
		StartYear:      startYear,
		EndYear:        endYear,
		RuntimeMinutes: runtimeMinutes,
		Genres:         genres,
	}
}

func InsertDatabaseActors() {

	path := "data/name.basics.tsv"

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	handleError(err)
	defer file.Close()

	fmt.Println("Inserting actors into database")

	transaction, err := db.Begin()
	handleError(err)

	stmt, err := transaction.Prepare(pq.CopyIn("actors", "nconst", "primaryname", "birthyear", "deathyear", "primaryprofession", "knownfortitles"))
	handleError(err)

	reader := bufio.NewReader(file)
	i := 0
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if i > 0 {
			actor := parseActor(string(line))
			_, err = stmt.Exec(actor.Nconst, actor.PrimaryName, actor.BirthYear, actor.DeathYear, pq.Array(actor.PrimaryProfession), pq.Array(actor.KnownForTitles))
			handleError(err)
		}
		i++
	}

	_, err = stmt.Exec()
	handleError(err)

	err = stmt.Close()
	handleError(err)

	err = transaction.Commit()
	handleError(err)

	fmt.Println("Finished inserting actors into database")

}

func parseActor(line string) Actor {
	arr := strings.Split(line, "\t")
	nconst := arr[0]
	primaryName := arr[1]
	birthYear, _ := strconv.Atoi(arr[2])
	deathYear, _ := strconv.Atoi(arr[3])
	primaryProfession := strings.Split(arr[4], ",")
	knownForTitles := strings.Split(arr[5], ",")

	if primaryProfession[0] == "\\N" {
		primaryProfession = []string{"No Primary Profession"}
	}

	if knownForTitles[0] == "\\N" {
		knownForTitles = []string{"Titles Not Found"}
	}

	return Actor{
		Nconst:            nconst,
		PrimaryName:       primaryName,
		BirthYear:         birthYear,
		DeathYear:         deathYear,
		PrimaryProfession: primaryProfession,
		KnownForTitles:    knownForTitles,
	}
}

func insertDatabaseMoviesActors() {
	path := "data/title.principals.tsv"

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	handleError(err)
	defer file.Close()

	// remove foreign key constraints
	_, err = db.Exec("ALTER TABLE title_actors DROP CONSTRAINT title_actors_tconst_fkey")
	handleError(err)
	_, err = db.Exec("ALTER TABLE title_actors DROP CONSTRAINT title_actors_nconst_fkey")
	handleError(err)

	fmt.Println("Inserting data into title_actors table...")

	transaction, err := db.Begin()
	handleError(err)

	stmt, err := transaction.Prepare(pq.CopyIn("title_actors", "tconst", "nconst", "title_ordering", "category", "job", "characters"))
	handleError(err)

	reader := bufio.NewReader(file)
	i := 0
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if i > 0 {
			movieActor := parseTitleActor(string(line))
			_, err = stmt.Exec(movieActor.Tconst, movieActor.Nconst, movieActor.Title_ordering, movieActor.Category, movieActor.Job, movieActor.Characters)
			handleError(err)
		}
		i++
	}

	_, err = stmt.Exec()
	handleError(err)

	err = stmt.Close()
	handleError(err)

	err = transaction.Commit()
	handleError(err)

	// remove all rows that don't have a corresponding title or actor
	_, err = db.Exec("delete from title_actors where tconst in (select title_actors.tconst from title_actors left join titles on title_actors.tconst = titles.tconst where titles.tconst is null)")
	handleError(err)

	_, err = db.Exec("delete from title_actors where nconst in (select title_actors.nconst from title_actors left join actors on title_actors.nconst = actors.nconst where actors.nconst is null)")
	handleError(err)

	// regenerate foreign key constraints
	_, err = db.Exec("ALTER TABLE title_actors ADD CONSTRAINT title_actors_tconst_fkey FOREIGN KEY (tconst) REFERENCES titles(tconst) ON DELETE CASCADE")
	handleError(err)
	_, err = db.Exec("ALTER TABLE title_actors ADD CONSTRAINT title_actors_nconst_fkey FOREIGN KEY (nconst) REFERENCES actors(nconst) ON DELETE CASCADE")
	handleError(err)

	// add indexes for tconst and nconst
	_, err = db.Exec("CREATE INDEX title_actors_tconst_idx ON title_actors (tconst)")
	handleError(err)
	_, err = db.Exec("CREATE INDEX title_actors_nconst_idx ON title_actors (nconst)")
	handleError(err)

	fmt.Println("Done inserting data into title_actors table...")

}

func parseTitleActor(line string) MovieActor {
	arr := strings.Split(line, "\t")
	tconst := arr[0]
	nconst := arr[2]
	category := arr[3]
	job := arr[4]
	if job == "\\N" {
		job = "No Job"
	}
	character := arr[5]
	if character == "\\N" {
		character = "No Character"
	}
	order, _ := strconv.Atoi(arr[1])

	return MovieActor{
		Tconst:         tconst,
		Title_ordering: order,
		Nconst:         nconst,
		Category:       category,
		Job:            job,
		Characters:     character,
	}
}

func ProcessItem(DItem DonwloadItem, wg *sync.WaitGroup) {
	Download(DItem)
	ExtractItem(DItem)
	defer wg.Done()
}

func ExtractItem(DItem DonwloadItem) {

	file, err := os.Open(DItem.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gzr.Close()

	// create new file
	newFile, err := os.Create(DItem.filename[:len(DItem.filename)-3])
	if err != nil {
		log.Fatal(err)
	}

	// copy data from gzr to newFile
	_, err = io.Copy(newFile, gzr)
	if err != nil {
		log.Fatal(err)
	}
	newFile.Close()

	// remove compressed file
	defer os.Remove(DItem.filename)
}

func Download(DItem DonwloadItem) {
	fmt.Println("Downloading: ", DItem.url)
	filename := DItem.filename
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		log.Fatal(err)
	}

	// Get the data
	client := http.Client{}
	resp, err := client.Get(DItem.url)
	if err != nil {
		fmt.Println("Error getting data: ", err)
		log.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		file.Write(line)
		if err == io.EOF {
			break
		}
	}

	defer resp.Body.Close()
	file.Close()
	fmt.Println("Downloaded: ", DItem.url)
}
