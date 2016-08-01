package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

type SampleResponse struct {
	Id    string    `json:"id"`
	SubID string    `json:"subID"`
	Time  time.Time `json:"time"`
}

func processResponse(resp *http.Response) {
	defer resp.Body.Close()

	var sampleResponse SampleResponse
	err := json.NewDecoder(resp.Body).Decode(&sampleResponse)
	if err != nil {
		logger.Printf("json err %v", err)
	} else {
		logger.Printf("sampleResponse %+v", sampleResponse)
	}
}

type urlRequest struct {
	url       string
	waitGroup *sync.WaitGroup
}

func processRequest(taskID int, request *urlRequest) {
	defer request.waitGroup.Done()

	logger.Printf("%v calling %v", taskID, request.url)
	resp, err := http.Get(request.url)
	if err != nil {
		logger.Printf("http get err %v", err)
	} else {
		processResponse(resp)
	}
}

func urlCallTask(taskID int, urlRequestChannel chan *urlRequest) {
	for request := range urlRequestChannel {
		processRequest(taskID, request)
	}
}

func main() {
	urlRequestChannel := make(chan *urlRequest)
	var waitGroup sync.WaitGroup
	for i := 0; i < 10; i += 1 {
		go urlCallTask(i, urlRequestChannel)
	}
	for i := 0; i < 1000; i += 1 {
		url := "http://localhost:8080/test/v1/1234/sub/" + strconv.Itoa(i)
		request := urlRequest{
			url:       url,
			waitGroup: &waitGroup,
		}
		waitGroup.Add(1)
		urlRequestChannel <- &request
	}
	waitGroup.Wait()
}
