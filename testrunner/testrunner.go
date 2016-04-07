package testrunner

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"log"
)

func RunTestContainers(client dockerclient.BurnsDockerClient) error {
	//Get images with test label and create containers for them
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false, Filters:map[string][]string{"label": {"test="}}})
	log.Printf("images: %+v", imgs)
	for i, img := range imgs {
		resultsPath := img.Labels[dockerclient.TestResultsLabel]
		c := dockerclient.NewContainer(&docker.Container{
			Name:        fmt.Sprintf("simpsons-%d", i),
			Config:     &docker.Config{ Image: img.ID },
			HostConfig: &docker.HostConfig{ Binds: []string{fmt.Sprintf("%s:%s", resultsPath, "/tmp")}},
		}, &docker.Image{ID:img.ID, }, )
		err := client.StartContainer(*c)
		if err != nil {
			log.Fatal("Failed starting container. ", err)
			return err
		}
	}

	return nil
}


