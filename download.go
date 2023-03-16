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

	"github.com/blockloop/scan"
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

	transaction, err := db.Begin()
	handleError(err)

	stmt, err := transaction.Prepare(pq.CopyIn("titles", "tconst", "titletype", "primarytitle", "originaltitle", "isadult", "startyear", "endyear", "runtimeminutes", "genres"))
	handleError(err)

	reader := bufio.NewReader(file)
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		title := parseTitle(string(line))
		_, err = stmt.Exec(title.Tconst, title.TitleType, title.PrimaryTitle, title.OriginalTitle, title.IsAdult, title.StartYear, title.EndYear, title.RuntimeMinutes, pq.Array(title.Genres))
		handleError(err)
	}
	_, err = stmt.Exec()
	handleError(err)

	err = stmt.Close()
	handleError(err)

	err = transaction.Commit()
	handleError(err)

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

	transaction, err := db.Begin()
	handleError(err)

	stmt, err := transaction.Prepare(pq.CopyIn("actors", "nconst", "primaryname", "birthyear", "deathyear", "primaryprofession", "knownfortitles"))
	handleError(err)

	reader := bufio.NewReader(file)
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		actor := parseActor(string(line))
		_, err = stmt.Exec(actor.Nconst, actor.PrimaryName, actor.BirthYear, actor.DeathYear, pq.Array(actor.PrimaryProfession), pq.Array(actor.KnownForTitles))
		handleError(err)
	}

	_, err = stmt.Exec()
	handleError(err)

	err = stmt.Close()
	handleError(err)

	err = transaction.Commit()
	handleError(err)

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

	transaction, err := db.Begin()
	handleError(err)

	stmt, err := transaction.Prepare(pq.CopyIn("title_actors", "tconst", "nconst", "title_ordering", "category", "job", "characters"))
	handleError(err)

	var missingTitles = make(map[string]int)
	var missingTitlesList []string

	reader := bufio.NewReader(file)
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		movieActor := parseTitleActor(string(line))

		if titleExists(movieActor.Tconst) {
			_, err = stmt.Exec(movieActor.Tconst, movieActor.Nconst, movieActor.Ordering, movieActor.Category, movieActor.Job, movieActor.Characters)
			handleError(err)
		} else {
			missingTitles[movieActor.Tconst] = 1
			_, contains := missingTitles[movieActor.Tconst]
			if contains {
				missingTitlesList = append(missingTitlesList, movieActor.Tconst)
			}
			fmt.Println("Missing Title: ", movieActor.Tconst)
		}
	}

	fmt.Println("Missing Titles: ", missingTitlesList)
	fmt.Println("Missing Titles Count: ", len(missingTitlesList))
	fmt.Println("Missing Titles: ", missingTitles)

	_, err = stmt.Exec()
	handleError(err)

	err = stmt.Close()
	handleError(err)

	err = transaction.Commit()
	handleError(err)
}

func titleExists(tconst string) bool {
	rows, err := db.Query("SELECT count(1) FROM titles WHERE tconst = $1 LIMIT 1", tconst)
	handleError(err)
	var affected int
	err = scan.Row(&affected, rows)
	return affected == 1
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
		Tconst:     tconst,
		Ordering:   order,
		Nconst:     nconst,
		Category:   category,
		Job:        job,
		Characters: character,
	}
}

func ProcessItem(DItem DonwloadItem, wg *sync.WaitGroup) {
	Download(DItem)
	ExtractItem(DItem)
	trimFirstLine(DItem.filename)
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

func trimFirstLine(filename string) {
	// remove the first line of each file
	// because it contains the column names
	filename = filename[:len(filename)-3]
	file, err := os.Open(filename)
	handleError(err)

	newFile, err := os.Create(filename + ".new")
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
		if i != 0 {
			newFile.Write(line)
			newFile.Write([]byte("\n"))
		}
		i++
	}

	// remove the file with extra line at start
	os.Remove(filename)

	// rename the file with the extra line removed
	os.Rename(filename+".new", filename)
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
