package task

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"os"
	"log"
	"errors"
	"path/filepath"
)

type DockerManager struct {
	client dockerclient.DockerClient
}

func NewDockerManager(dockerClient dockerclient.DockerClient) DockerManager {

	return DockerManager{client: dockerClient, }
}

func (manager DockerManager) GetImages() ([]docker.APIImages, error) {

	return manager.client.ListImages(
		docker.ListImagesOptions{All: false,
			Filters: map[string][]string{"label": {"test="}}});
}

func (manager DockerManager) RunTests(image docker.APIImages, containerName string) (string, error) {

	containerResultsPath := image.Labels[dockerclient.LabelTestResultPath]
	containerCmd := image.Labels[dockerclient.LabelTestCmd]
	resultDirName := fmt.Sprintf("/tmp/test-results/%s", containerName)
	os.MkdirAll(resultDirName, 0700)
	containerConfig := &docker.Config{Image: image.ID}
	if len(containerCmd) > 0 {
		containerConfig.Cmd = []string{containerCmd }
	}
	err := manager.startContainer(image, containerName, containerConfig,
		&docker.HostConfig{Binds: []string{fmt.Sprintf("%s:%s", resultDirName, containerResultsPath)}})
	if err != nil {
		log.Printf("Failed starting container: %s. Error: %v", containerName, err)
		return "", err
	}
	status, err := manager.client.WaitContainer(containerName)
	if err != nil {
		log.Printf("Failed waiting for container: %s. Error: %v", containerName, err)
		return "", err
	}
	if status != 0 {
		return "", errors.New(fmt.Sprintf("Failed waiting for container: %s. Status: %v", containerName, status))
	}

	return filepath.Join(resultDirName, image.Labels[dockerclient.LabelTestResultsFile]), nil
}

func (manager DockerManager) GetImageRunningInterval(image docker.APIImages) string {

	return image.Labels[dockerclient.LabelInterval]
}

func (manager DockerManager) startContainer(image docker.APIImages, containerName string, containerConfig *docker.Config, hostConfig *docker.HostConfig) error {

	manager.client.RemoveContainer(containerName, true)
	c := dockerclient.NewContainer(&docker.Container{
		Name:        containerName,
		Config:     containerConfig,
		HostConfig: hostConfig,
	}, &docker.Image{ID: image.ID})
	err := manager.client.StartContainer(*c)

	return err
}
