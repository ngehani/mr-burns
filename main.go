package main

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/mr-burns/controller"
	"github.com/gaia-adm/mr-burns/common"
)

func main() {

	endpoint := "unix:///var/run/docker.sock"
	client := controller.NewDockerManager(dockerclient.NewClient(endpoint))
	publisher := controller.NewPublisher(common.NewConfiguration().PublisherURL)
	controller.NewController(client, publisher).Start()
}