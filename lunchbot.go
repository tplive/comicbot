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
	"time"

	"github.com/joho/godotenv"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

func main() {

	webHookUrl, ok := getEnvironment("WEBHOOK_URL")
	if !ok {
		log.Fatal("No string")
	}

	fileName, url := getComic()

	err1 := SendSlackNotification(webHookUrl, "Dagens Lunch "+url)
	if err1 != nil {
		log.Fatal(err1)
	}

	err2 := downloadFile(url, fileName)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Printf("File %s downloaded in current working directory", fileName)
}

func getEnvironment(varName string) (string, bool) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.LookupEnv(varName)
}

func HowToSendSlackNotification() {
	webHookUrl, _ := getEnvironment("WEBHOOK_URL")
	_, url := getComic()
	err := SendSlackNotification(webHookUrl, "Dagens Lunch "+url)
	if err != nil {
		log.Fatal(err)
	}
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
		return errors.New("Non-ok repsonse returned from Slack")
	}
	return nil
}

func getComic() (string, string) {
	baseUrl := "https://tu.no"
	date := time.Now().Format("2006-01-02")
	comicId := "lunch"
	fileName := "tu-" + comicId + "-" + date + ".jpg"
	URL := baseUrl + "/?module=TekComics&service=image&id=" + comicId + "&key=" + date
	// URL looks like this https://www.tu.no/?module=TekComics&service=image&id=lunch&key=2020-05-28

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
