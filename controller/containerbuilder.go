package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"fmt"
	"os"
	"encoding/json"
)

func BuildContainer(image docker.APIImages, containerName string, resultDirName string) dockerclient.Container {

	ret := createContainerSettings(image)
	ret.Name = containerName
	bindResultDir(&ret, image, containerName)
	bindImageId(&ret, image)

	return dockerclient.NewContainer(&ret)
}

func createContainerSettings(image docker.APIImages) docker.Container {

	containerSettingsJson := getContainerSettingsJson(image)
	var ret docker.Container
	err := json.Unmarshal([]byte(containerSettingsJson), &ret)
	if err != nil {
		log.Errorf("Failed to create container from 'test.container.settings' label. Check JSON format, for detailed format look at https://docs.docker.com/reference/api/docker_remote_api_v1.16/#create-a-container (%s)", containerSettingsJson)
		panic(err)
	}

	return ret
}

func getContainerSettingsJson(image docker.APIImages) string {

	ret := image.Labels[dockerclient.LABEL_TEST_CONTAINER_SETTINGS]
	if ret == "" {
		ret = "{}"
	}

	return ret
}

func bindResultDir(containerSettings *docker.Container, image docker.APIImages, resultDirName string) {

	containerResultsPath := image.Labels[dockerclient.LABEL_TEST_RESULTS_DIR]
	os.MkdirAll(resultDirName, 0700)
	if (containerSettings.HostConfig == nil) {
		containerSettings.HostConfig = &docker.HostConfig{}
	}
	if (containerSettings.HostConfig.Binds == nil) {
		containerSettings.HostConfig.Binds = []string{}
	}
	containerSettings.HostConfig.Binds = append(containerSettings.HostConfig.Binds, fmt.Sprintf("%s:%s", resultDirName, containerResultsPath))
}

func bindImageId(containerSettings *docker.Container, image docker.APIImages) {

	if (containerSettings.Config == nil) {
		containerSettings.Config = &docker.Config{}
	}
	containerSettings.Config.Image = image.ID
}