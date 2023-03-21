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

	// Pseudo code:

	// 1. Get info from latest comic.
	// 2. Get the index of the last downloaded comic from persistence layer.
	// 3. If there is a new comic, get it.
	// 4. Post to Slack
	// 5. Download image file
	// 6. If successful, increment index in persistence layer.

	data := getLatestComicMetadata()
	lastComic := getLastComicIndex(counterApiUrl, counterApiBucket)

	if data.Num > lastComic {
		// There are newer comic(s), so we should get'em
		fmt.Printf("Url for the \"next\" comic: %s\n", data.Img)
		fileName := fmt.Sprintf("xkcd-%d.png", data.Num)
		err := downloadFileFromUrl(data.Img, fileName)
		if err == nil {
			incrementComicIndex(data.Num)
		}

	}
	fmt.Printf("Current comic: %d\n", data.Num)
	fmt.Printf("Last fetched comic: %d\n", lastComic)
}

func incrementComicIndex(i int) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(json.Marshal("+1")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

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

func getLatestComicMetadata() XKCDComic {

	// Current comic: https://xkcd.com/info.0.json
	// Arbitrary comic: https://xkcd.com/{comicIndex}/info.0.json where comicIndex is 1..N
	currentComicUrl := "https://xkcd.com/info.0.json"

	response, apiError := http.Get(currentComicUrl)

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

	//fmt.Printf("Results: %v\n", comicMetadata)
	//fmt.Println(comicMetadata.Num)
	//fmt.Println(comicMetadata.SafeTitle)

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
