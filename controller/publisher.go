package controller

import (
	"bytes"
	"net/http"
	log "github.com/Sirupsen/logrus"
)

type Publisher struct {
	url string
}

func NewPublisher(webhookUrl string) Publisher {

	return Publisher{url: webhookUrl}
}

func (publisher Publisher) Publish(data string) error {

	req, err := http.NewRequest("POST", publisher.url, bytes.NewBufferString(data))
	if err != nil {
		log.Error("Failed to create HTTP request", data, err)
		return err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Error("Failed to POST container test results", data, err)
		return err
	}
	defer response.Body.Close()

	return nil
}
