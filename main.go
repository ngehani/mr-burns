package main

import (
	"github.com/gaia-adm/mr-burns/controller"
	"github.com/gaia-adm/mr-burns/dockerclient"
)

func main() {

	endpoint := "unix:///var/run/docker.sock"
	client := dockerclient.NewClient(endpoint)
	controller.Start(client)
}
