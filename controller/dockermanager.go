package controller

import (
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"os"
	log "github.com/Sirupsen/logrus"
	"errors"
	"path/filepath"
	"strings"
	"github.com/gaia-adm/mr-burns/common"
)

type DockerManager struct {
	client dockerclient.DockerClient
}

func NewDockerManager(dockerClient dockerclient.DockerClient) DockerManager {

	return DockerManager{client: dockerClient, }
}

func (manager DockerManager) GetImages() ([]docker.APIImages, error) {

	images, err := manager.client.ListImages(
		docker.ListImagesOptions{All: false, Filters: map[string][]string{"label": {"test="}, "dangling": {"false"}}})
	if err != nil {
		return nil, err
	}
	log.Infof("Fetched images: %+v", dangling(images))

	// filter dangling (doesn't supported as part of docker API like on docker swarm)
	return dangling(images), nil
}

func (manager DockerManager) RunTests(image docker.APIImages, containerName string) (string, error) {

	containerResultsPath := image.Labels[dockerclient.LABEL_TEST_RESULTS_DIR]
	containerCmd := image.Labels[dockerclient.LABEL_TEST_CMD]
	resultsDirName := fmt.Sprintf("/tmp/test-results/%s_%d", containerName, common.GetTimeNowMillisecond())
	os.MkdirAll(resultsDirName, 0700)
	containerConfig := &docker.Config{Image: image.ID}
	if len(containerCmd) > 0 {
		containerConfig.Cmd = []string{containerCmd }
	}
	err := manager.startContainer(image, containerName, containerConfig,
		&docker.HostConfig{Binds: []string{fmt.Sprintf("%s:%s", resultsDirName, containerResultsPath)}})
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

	return filepath.Join(resultsDirName, image.Labels[dockerclient.LABEL_TEST_RESULTS_FILE]), nil
}

func (manager DockerManager) startContainer(image docker.APIImages, containerName string, containerConfig *docker.Config, hostConfig *docker.HostConfig) error {

	manager.client.RemoveContainer(containerName, true)
	c := dockerclient.NewContainer(&docker.Container{
		Name:        containerName,
		Config:     containerConfig,
		HostConfig: hostConfig,
	}, &docker.Image{ID: image.ID})

	return manager.client.StartContainer(*c)
}

func (manager DockerManager) GetLabelImageRunningInterval(image docker.APIImages) string {

	return image.Labels[dockerclient.LABEL_INTERVAL]
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