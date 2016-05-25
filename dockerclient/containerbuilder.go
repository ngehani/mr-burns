package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"encoding/json"
)

func BuildContainer(image docker.APIImages, containerName string, resultDirName string) Container {

	ret := createContainerSettings(image)
	ret.Name = containerName
	bindResultDir(&ret, image, containerName)
	bindImageId(&ret, image)

	return NewContainer(&ret)
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

	ret := ContainerSettings(image)
	if ret == "" {
		ret = "{}"
	}

	return ret
}

func bindResultDir(containerSettings *docker.Container, image docker.APIImages, resultDirName string) {

	containerResultsPath := ResultsDir(image)
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