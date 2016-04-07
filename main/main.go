package main

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/mr-burns/testrunner"
)

func main() {

	endpoint := "http://gaia-local.skydns.local:2375"
	client := dockerclient.NewClient(endpoint, true)
	testrunner.RunTestContainers(client)
}