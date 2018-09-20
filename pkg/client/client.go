package client

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

var (
	timeout = time.Second * 3

	client = &http.Client{Timeout: timeout}
)

//DoPOST uses for making http POST request to a specific URL
//client has set timeout for 3 seconds
func DoPOST(URL string, body []byte, logger *log.Logger) (*http.Response, error) {
	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Println("Couldn't send a request to server.", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Printf("Sent message, but listener gave respond with status code: [%d]", resp.StatusCode)
	}

	return resp, err
}
