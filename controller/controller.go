package controller

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/go-dockerclient"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"bytes"
	"net/http"
)

func Start(client dockerclient.DockerClient) {

	// Get images with test label and create containers for them
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false, Filters:map[string][]string{"label": {"test="}}})
	log.Printf("images: %+v", imgs)
	for i, img := range imgs {
		runTestContainer(client, img, fmt.Sprintf("simpsons-%d", i))
	}
}

func runTestContainer(client dockerclient.DockerClient, image docker.APIImages, containerName string) error {

	containerResultsPath := image.Labels[dockerclient.LabelTestResultPath]
	containerResultsFile := image.Labels[dockerclient.TestResultsFileLabel]
	containerCmd := image.Labels[dockerclient.TestCmdLabel]
	resultDirName := fmt.Sprintf("/tmp/test-results/%s", containerName)
	os.MkdirAll(resultDirName, 0700)
	containerConfig := &docker.Config{Image: image.ID }
	if len(containerCmd) > 0 {
		containerConfig.Cmd = []string{containerCmd }
	}
	err := startContainer(client,
		containerName,
		containerConfig,
		&docker.HostConfig{Binds: []string{fmt.Sprintf("%s:%s", resultDirName, containerResultsPath)}},
		image)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed starting container: %s.", containerName), err)
		return err
	}

	status, err := client.WaitContainer(containerName)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed waiting for container: %s.", containerName), err)
		return err
	}
	if (status == 0) {
		postResults(resultDirName, containerResultsFile)
	}

	return nil
}

func startContainer(client dockerclient.DockerClient, containerName string, containerConfig *docker.Config, hostConfig *docker.HostConfig, image docker.APIImages) error {

	client.RemoveContainer(containerName, true)
	c := dockerclient.NewContainer(&docker.Container{
		Name:        containerName,
		Config:     containerConfig,
		HostConfig: hostConfig,
	}, &docker.Image{ID:image.ID, }, )
	err := client.StartContainer(*c)

	return err
}

func postResults(resultDirName string, containerResultsFile string) error {

	containerTestResults, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", resultDirName, containerResultsFile))
	if (err != nil) {
		log.Fatal("Failed to read file", err)
		return err
	}
	req, err := http.NewRequest("POST", "http://distributor-link:8000", bytes.NewBuffer(containerTestResults))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to POST container test results", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}


