package main

import (
	"fmt"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
)

func main() {

	endpoint := "unix:///var/run/docker.sock"
	client := dockerclient.NewClient(endpoint, true)
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false, Filters:map[string][]string{"label": {"test="}}})
	for _, img := range imgs {
		fmt.Println("img: %+v", img)
	}
}