package testrunner

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"bytes"
	"net/http"
)

func RunTestContainers(client dockerclient.BurnsDockerClient) {

	//Get images with test label and create containers for them
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false, Filters:map[string][]string{"label": {"test="}}})
	log.Printf("images: %+v", imgs)
	for i, img := range imgs {
		runTestContainer(client, img, fmt.Sprintf("simpsons-%d", i))
	}
}

func runTestContainer(client dockerclient.BurnsDockerClient, image docker.APIImages, containerName string) error {

	resultsPath := image.Labels[dockerclient.TestResultsLabel]
	resultDirName := fmt.Sprintf("/tmp/test-results/%s", containerName)
	os.MkdirAll(resultDirName, 0700)
	client.RemoveContainer(containerName, true)
	c := dockerclient.NewContainer(&docker.Container{
		Name:        containerName,
		Config:     &docker.Config{ Image: image.ID },
		HostConfig: &docker.HostConfig{ Binds: []string{fmt.Sprintf("%s:%s", resultDirName, resultsPath)}},
	}, &docker.Image{ID:image.ID, }, )
	err := client.StartContainer(*c)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed starting container: %s.", containerName), err)
		return err
	}
	status, err := client.WaitContainer(containerName)
	if (status == 0) {
		containerTestResults, err := ioutil.ReadFile(fmt.Sprintf("%s/TestSuite.txt", resultDirName))
		if(err != nil) {
			log.Fatal("Failed to read file", err)
		}
		req, err := http.NewRequest("POST", "http://distributor-link:8000", bytes.NewBuffer(containerTestResults))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Failed to POST container test results", err)
		}
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed waiting for container: %s.", containerName), err)
		return err
	}

	return nil
}


