package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const url string = "https://xkcd.com"
const counterApiUrl string = "https://kvdb.io"
const counterApiBucket string = "32FZdnyZxxV3Rm7VWqRWH4"
const endpoint string = counterApiUrl + "/" + counterApiBucket + "/xkcd"

type XKCDComic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

func getXKCD() {

	// Get metadata for the latest posted XKCD comic strip. -1 means "get latest"
	latestComicMetadata := getComicMetadata(-1)

	// Get the index of the last comic we downloaded.
	previousComic := getLastComicIndex(counterApiUrl, counterApiBucket)
	var lastComic int
	lastComic = previousComic

	// Get Slack URL for posting
	webHookUrl := getEnvVar("WEBHOOK_URL")
	if webHookUrl == "" {
		log.Fatal("No such environment variable WEBHOOK_URL")
	}

	// Sort of a while loop - if there are newer comic(s) since last update, we should get'em all
	for latestComicMetadata.Num > previousComic {

		// We need to know the last comic we downloaded
		nextComic := getComicMetadata(previousComic)
		var err error

		sendSlackNotification(webHookUrl, "Siste XKCD "+nextComic.Img)

		fileName := fmt.Sprintf("xkcd-%d.png", nextComic.Num)
		err = downloadFileFromUrl(nextComic.Img, fileName)
		if err == nil {
			previousComic++
		} else {
			break
		}

		// Break out of the loop once we have downloaded the latest comic
		if latestComicMetadata.Num == previousComic {
			setComicIndex(latestComicMetadata.Num)
			break
		}
	}
}

func incrementComicIndex() {
	client := &http.Client{}

	payload := []byte("+1")

	req, _ := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func setComicIndex(value int) {
	client := &http.Client{}

	payload := []byte(fmt.Sprint(value))

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func getLastComicIndex(url string, bucket string) int {

	// GET kvdb.io/bucket/key
	response, apiError := http.Get(counterApiUrl + "/" + bucket + "/xkcd")

	if apiError != nil {
		log.Fatal(apiError)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var latestComic int
	err = json.Unmarshal(body, &latestComic)

	if err != nil {
		panic(err.Error())
	}

	return latestComic
}

func getComicMetadata(comicNumber int) XKCDComic {

	// Arbitrary comic: https://xkcd.com/{comicIndex}/info.0.json where comicIndex is 1..N
	// Latest published comic: https://xkcd.com/info.0.json
	// If comicNumber is set to a positive integer, it means we want one particular comic. If not, get the latest one.

	var comicUrl string
	if comicNumber > 0 {
		comicUrl = fmt.Sprintf("https://xkcd.com/%d/info.0.json", comicNumber)
	} else {
		comicUrl = "https://xkcd.com/info.0.json"
	}

	response, apiError := http.Get(comicUrl)

	if apiError != nil {
		log.Fatal(apiError)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var comicMetadata XKCDComic

	err = json.Unmarshal(body, &comicMetadata)
	if err != nil {
		panic(err.Error())
	}

	return comicMetadata
}

func XKCDCreate(w http.ResponseWriter, r *http.Request) {

	var x XKCDComic

	err := decodeJSONBody(w, r, &x)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Comic: %+v", x)
}
