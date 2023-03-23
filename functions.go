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
	"strings"
	"time"
	"unicode"

	"github.com/joho/godotenv"
)

// Common functions for Comicbot

func downloadFileFromUrl(URL, fileName string) error {
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

func getEnvVar(varName string) string {

	godotenv.Load(".env")

	value, isSet := os.LookupEnv(varName)

	if !isSet || value == "" {
		log.Print("Must set environment variable " + varName)
	}
	return value
}

func getLastComic(comic string, apiUrl string, apiBucket string) int16 {
	// Get the previously stored comic
	response, err := http.Get("$apiUrl/$apiBucket/$comic")
	if err != nil {
		log.Fatal("Unable to access " + apiUrl + "/" + apiBucket + " - does it exist?")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal("Received non 200 response code: " + strconv.Itoa(response.StatusCode))

	}

	return 1

}

func sendSlackNotification(webhookUrl string, msg string) error {

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

func capitalize(s string) string {
	sep := " "
	ss := strings.SplitN(s, sep, 2)
	r := []rune(ss[0])
	if len(r) == 0 {
		return s
	}
	r[0] = unicode.ToUpper(r[0])
	s = string(r)
	if len(ss) > 1 {
		s += sep + ss[1]
	}
	return s
}
