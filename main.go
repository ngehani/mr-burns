package main

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/mr-burns/testrunner"
)

func main() {

	endpoint := "unix:///var/run/docker.sock"
	client := dockerclient.NewClient(endpoint, true)
	testrunner.RunTestContainers(client)
}