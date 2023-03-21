package main

import (
	"log"
	"time"
)

const baseUrl string = "https://tu.no/api/widgets/comics?name="
const oldBaseUrl string = "https://tu.no/modules/?module=TekComics&service=image&id="

type SlackRequestBody struct {
	Text string `json:"text"`
}

func main() {

	// Setup
	webHookUrl := getEnvVar("WEBHOOK_URL")
	if webHookUrl == "" {
		log.Fatal("No such environment variable WEBHOOK_URL")
	}

	// List comics from Teknisk Ukeblad
	tuComix := []string{
		"lunch",
		"dilbert",
		"dunce",
	}

	// Iterate over the comics
	for _, comic := range tuComix {
		fileName, url := getTekniskUkebladComic(comic)

		// Post comic to Slack
		notErr := sendSlackNotification(webHookUrl, "Dagens "+toUpper(comic)+" "+url)
		if notErr != nil {
			log.Fatal(notErr)
		}

		// Download image
		fileErr := downloadFileFromUrl(url, fileName)
		if fileErr != nil {
			log.Fatal(fileErr)
		}

	}
}

func getTekniskUkebladComic(comic string) (string, string) {
	date := time.Now().Format("2006-01-02")

	var URL string
	// The old base URL looks like this https://www.tu.no/modules/?module=TekComics&service=image&id=dilbert&key=2020-05-28 and is still used for Dilbert
	// The new format introduced in March 2023 looks like this https://www.tu.no/api/widgets/comics?name=lunch&date=2023-03-14 and is used for Lunch and Dunce

	if comic == "dilbert" {

		URL = oldBaseUrl + comic + "&key=" + date
	} else if comic == "lunch" || comic == "dunce" {
		URL = baseUrl + comic + "&date=" + date
	} else {
		URL = oldBaseUrl + comic + "&key=" + date
	}

	fileName := "tu-" + comic + "-" + date + ".jpg"

	println(URL)

	return fileName, URL
}
