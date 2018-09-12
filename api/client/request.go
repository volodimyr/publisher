package client

import (
	"bytes"
	"github.com/volodimyr/event_publisher/api/config"
	"log"
	"net/http"
)

func DoPOST(URL string, body []byte) {
	resp, err := http.Post(URL, config.ContentType, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Couldn't send a request to server.", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("Sent message, but listener gave respond with status code: [%d]", resp.StatusCode)
	}
}
