package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"fmt"
	"github.com/gaia-adm/mr-burns/common"
	"os"
	"encoding/json"
)

func BuildContainer(image docker.APIImages, containerName string) docker.Container {

	ret := createContainerSettings(image)
	bindResultDir(&ret, image, containerName)
	bindImageId(&ret, image)

	return ret
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

	return image.Labels[dockerclient.LABEL_TEST_CONTAINER_SETTINGS]
}

func bindResultDir(containerSettings *docker.Container, image docker.APIImages, containerName string) {

	containerResultsPath := image.Labels[dockerclient.LABEL_TEST_RESULTS_DIR]
	resultsDirName := fmt.Sprintf("/tmp/test-results/%s_%d", containerName, common.GetTimeNowMillisecond())
	os.MkdirAll(resultsDirName, 0700)
	if (containerSettings.HostConfig == nil) {
		containerSettings.HostConfig = &docker.HostConfig{}
	}
	if (containerSettings.HostConfig.Binds == nil) {
		containerSettings.HostConfig.Binds = []string{}
	}
	containerSettings.HostConfig.Binds = append(containerSettings.HostConfig.Binds, fmt.Sprintf("%s:%s", resultsDirName, containerResultsPath))
}

func bindImageId(containerSettings *docker.Container, image docker.APIImages) {

	if (containerSettings.Config == nil) {
		containerSettings.Config = &docker.Config{}
	}
	containerSettings.Config.Image = image.ID
}