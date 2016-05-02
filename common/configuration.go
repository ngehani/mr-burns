package common

//import (
//	"encoding/json"
//	"os"
//)

type Configuration struct {
	PublisherURL string
}

func NewConfiguration() Configuration {

	return Configuration{PublisherURL: "http://distributor.skydns.local"}

	//file, err := os.Open("../mr-burns-configuration.json")
	//if err != nil {
	//	panic(err)
	//}
	//decoder := json.NewDecoder(file)
	//ret := Configuration{}
	//err = decoder.Decode(&ret)
	//if err != nil {
	//	panic(err)
	//}
	//
	//return ret
}