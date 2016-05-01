package common

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	PublisherURL string
}

func NewConfiguration() Configuration {

	file, err := os.Open("../mr-burns-configuration.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	ret := Configuration{}
	err = decoder.Decode(&ret)
	if err != nil {
		panic(err)
	}

	return ret
}