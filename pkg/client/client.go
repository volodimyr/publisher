package client

import (
	"bytes"
	"log"
	"net/http"
)

//DoPOST uses for making http POST request to a specific URL
func DoPOST(URL string, body []byte, logger *log.Logger) {
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Println("Couldn't send a request to server.", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Printf("Sent message, but listener gave respond with status code: [%d]", resp.StatusCode)
	}
}
