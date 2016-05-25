package controller

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"errors"
	"path/filepath"
	"strings"
	"github.com/gaia-adm/mr-burns/common"
	"os"
)

type DockerManager struct {
	client dockerclient.DockerClient
}

func NewDockerManager(dockerClient dockerclient.DockerClient) DockerManager {

	return DockerManager{client: dockerClient, }
}

func (manager DockerManager) GetImages() ([]docker.APIImages, error) {

	images, err := manager.client.ListImages(
		docker.ListImagesOptions{All: false, Filters: map[string][]string{
			"label": {fmt.Sprintf("%s=true", dockerclient.LABEL_TEST)},
			"dangling": {"false"}}})
	if err != nil {
		return nil, err
	}
	log.Infof("Fetched images: %+v", dangling(images))

	// filter dangling (doesn't supported as part of docker API like on docker swarm)
	return dangling(images), nil
}

func (manager DockerManager) RunTests(image docker.APIImages, containerName string) (string, error) {

	resultDirName := createResultDir(containerName)
	container := dockerclient.BuildContainer(image, containerName, resultDirName)
	err := manager.startContainer(image, container)
	if err != nil {
		log.Infof("Failed starting container: %s. Error: %v", containerName, err)
		return "", err
	}
	status, err := manager.client.WaitContainer(containerName)
	if err != nil {
		log.Error("Failed waiting for container: %s. Error: %v", containerName, err)
		return "", err
	}
	log.Infof("Finish wating for container: %s, status: %d", containerName, status)
	if status != 0 {
		return "", errors.New(fmt.Sprintf("Failed waiting for container: %s. Status: %v", containerName, status))
	}

	return filepath.Join(resultDirName, dockerclient.ResultsFile(image)), nil
}

func (manager DockerManager) startContainer(image docker.APIImages, container dockerclient.Container) error {

	manager.client.RemoveContainer(container.Data.Name, true)

	return manager.client.StartContainer(container)
}

func (manager DockerManager) GetContainerLogs(container string) (string, error) {

	logs, err := manager.client.Logs(container)
	if err != nil {
		log.Error("Failed to get container logs. ", container, err)
	}

	return logs, err
}

func dangling(images []docker.APIImages) []docker.APIImages {

	var ret []docker.APIImages
	for _, currImage := range images {
		if !strings.Contains(currImage.RepoTags[0], "<none>") {
			ret = append(ret, currImage)
		}
	}

	return ret
}

func createResultDir(container string) string {

	ret := fmt.Sprintf("/tmp/test-results/%s_%d", container, common.GetTimeNowMillisecond())
	err := os.MkdirAll(ret, 0700)
	if err != nil {
		log.Errorf("Failed to create test results folder: %s. Error: %v", ret, err)
		panic(err)
	}

	return ret
}