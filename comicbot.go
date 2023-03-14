package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const baseUrl string = "https://tu.no/api/widgets"
const oldBaseUrl string = "https://tu.no/modules"

type SlackRequestBody struct {
	Text string `json:"text"`
}

func main() {

	webHookUrl := getEnvironment("WEBHOOK_URL")
	if webHookUrl == "" {
		log.Fatal("No string")
	}

	lunch, LUrl := getComic("lunch")
	dunce, DuUrl := getComic("dunce")
	dilbert, DUrl := getComic("dilbert")

	LNotErr := SendSlackNotification(webHookUrl, "Dagens Lunch "+LUrl)
	if LNotErr != nil {
		log.Fatal(LNotErr)
	}

	DNotErr := SendSlackNotification(webHookUrl, "Dagens Dilbert "+DUrl)
	if DNotErr != nil {
		log.Fatal(DNotErr)
	}

	DuNotErr := SendSlackNotification(webHookUrl, "Dagens Dunce "+DuUrl)
	if DuNotErr != nil {
		log.Fatal(DuNotErr)
	}

	LFileErr := downloadFile(LUrl, lunch)
	if LFileErr != nil {
		log.Fatal(LFileErr)
	}

	DuFileErr := downloadFile(DuUrl, dunce)
	if DuFileErr != nil {
		log.Fatal(DuFileErr)
	}

	DFileErr := downloadFile(DUrl, dilbert)
	if DFileErr != nil {
		log.Default()
	}

}

func getEnvironment(varName string) string {

	godotenv.Load(".env")

	value, isSet := os.LookupEnv(varName)

	if !isSet || value == "" {
		log.Print("Must set environment variable " + varName)
	}
	return value
}

func SendSlackNotification(webhookUrl string, msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("non-ok response returned from slack")
	}
	return nil
}

func getComic(comic string) (string, string) {
	date := time.Now().Format("2006-01-02")

	var URL string
	// The old base URL looks like this https://www.tu.no/modules/?module=TekComics&service=image&id=dilbert&key=2020-05-28 and is still used for Dilbert
	// The new format introduced in March 2023 looks like this https://www.tu.no/api/widgets/comics?name=lunch&date=2023-03-14 and is used for Lunch and Dunce

	if comic == "dilbert" {

		URL = oldBaseUrl + "/?module=TekComics&service=image&id=" + comic + "&key=" + date
	} else if comic == "lunch" || comic == "dunce" {
		URL = baseUrl + "/comics?name=" + comic + "&date=" + date
	} else {
		URL = oldBaseUrl + "/?module=TekComics&service=image&id=" + comic + "&key=" + date
	}

	fileName := "tu-" + comic + "-" + date + ".jpg"

	println(URL)

	return fileName, URL
}

func downloadFile(URL, fileName string) error {
	// Get response bytes from url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code: " + strconv.Itoa(response.StatusCode))

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
	fmt.Printf("File %s downloaded in current working directory\n", fileName)
	return nil
}
