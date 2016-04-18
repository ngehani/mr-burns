package main

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/mr-burns/controller"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	ls()
	endpoint := "unix:///var/run/docker.sock"
	client := dockerclient.NewClient(endpoint)
	controller.Start(client)
}

func ls() {

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.Getwd()
	log.Println("current directory: %s", dir)
	for _, file := range files {
		log.Printf("%+v", file)
	}
}