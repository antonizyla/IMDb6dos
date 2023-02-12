package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type DonwloadItem struct {
	url      string
	filename string
}

func DispatchDownloads() {

	os.Mkdir("data", 0777)
	var urls = [1]DonwloadItem{{"https://datasets.imdbws.com/title.ratings.tsv.gz", "data/title.ratings.tsv.gz"}}

	var workGroup sync.WaitGroup
	for _, url := range urls {
		workGroup.Add(1)
		go ProcessItem(url, &workGroup)
	}
	workGroup.Wait()
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

	gzr, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gzr.Close()
	defer file.Close()

	// create new file
	newFile, err := os.Create(DItem.filename[:len(DItem.filename)-3])
	if err != nil {
		log.Fatal(err)
	}

	// copy data from gzr to newFile
	scanner := bufio.NewScanner(gzr)
	for scanner.Scan() {
		newFile.Write(scanner.Bytes())
	}
	defer newFile.Close()

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