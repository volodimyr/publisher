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
//returns nil error if request was sent successfully
func DoPOST(URL string, body []byte, logger *log.Logger) (*http.Response, error) {
	resp, err := client.Post(URL, "application /json", bytes.NewBuffer(body))
	if err != nil {
		logger.Printf("Couldn't send a request [%s] to server: [%v]\n", string(body), err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Printf("Sent message [%s], but listener gave respond with status code: [%d]\n", string(body), resp.StatusCode)
	}

	return resp, err
}
