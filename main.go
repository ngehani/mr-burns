package main

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/mr-burns/controller"
)

func main() {

	endpoint := "unix:///var/run/docker.sock"
	client := dockerclient.NewClient(endpoint)
	controller.NewController(controller.NewDockerManager(client)).Start()
}
