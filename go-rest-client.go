package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

type SampleResponse struct {
	Id    string    `json:"id"`
	SubID string    `json:"subID"`
	Time  time.Time `json:"time"`
}

func processResponse(resp* http.Response) (err error) {
	defer resp.Body.Close()

	var sampleResponse SampleResponse
	err = json.NewDecoder(resp.Body).Decode(&sampleResponse)
	if err == nil {
		logger.Printf("sampleResponse %+v", sampleResponse)
	}
	return
}

func callURL(url string) (err error) {
	logger.Printf("calling %v", url)
	resp, err := http.Get(url)
	if err == nil {
		err = processResponse(resp)
	}
	return
}

func main() {
	url := "http://localhost:8080/test/v1/1234/sub/2345"
	if err := callURL(url); err != nil {
		logger.Printf("err %v", err)
	}
}
