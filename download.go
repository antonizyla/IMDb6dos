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
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// scanner couldn't handle large file or i did it wrong
	// table has 12 cols, 65k param limit for postgres => 5.45k rows per insert
	listTitles := [5400]Title{}
	i := 0
	reader := bufio.NewReader(file)
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if i >= 5400 {
			db.Create(&listTitles)
			i = 0
		}
		listTitles[i] = parseTitle(string(line))
		i++
	}

	reader = bufio.NewReader(file)

	db.Create(&listTitles)
	//fmt.Println("Lines: ", lines)

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

	genreField := strings.Split(arr[8], ",")
    genresList := make([]Genre, len(genreField))
	for i, genre := range genreField {
		genresList[i] = Genre{Genre: genre}
	}

	// simplest way to do in one pass
	// but it's not the best way to do

	//db.Create(&genresList)

	return Title{
		Tconst:         tconst,
		TitleType:      titleType,
		PrimaryTitle:   primaryTitle,
		OriginalTitle:  originalTitle,
		IsAdult:        isAdult,
		StartYear:      startYear,
		EndYear:        endYear,
		RuntimeMinutes: runtimeMinutes,
		Genres:         genresList,
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
