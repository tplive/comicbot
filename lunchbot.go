package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	baseUrl := "https://tu.no"
	date := time.Now().Format("2006-01-02")
	comicId := "lunch"
	fileName := "tu-" + comicId + "-" + date + ".jpg"
	URL := baseUrl + "/?module=TekComics&service=image&id=" + comicId + "&key=" + date
	// URL looks like this https://www.tu.no/?module=TekComics&service=image&id=lunch&key=2020-05-28

	err := downloadFile(URL, fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File %s downloaded in current working directory", fileName)

}

func downloadFile(URL, fileName string) error {
	// Get response bytes from url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	// Create an empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write bytes to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
