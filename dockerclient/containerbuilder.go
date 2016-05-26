package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"encoding/json"
	"github.com/gaia-adm/mr-burns/common"
)

func BuildContainer(image docker.APIImages, containerName string, resultDirName string) Container {

	ret := createContainerSettings(image)
	ret.Name = containerName
	bindEnv(&ret, image)
	bindResultDir(&ret, image, resultDirName)
	bindImageId(&ret, image)

	return NewContainer(&ret)
}

func createContainerSettings(image docker.APIImages) docker.Container {

	containerSettingsJson := getContainerSettingsJson(image)
	var ret docker.Container
	err := json.Unmarshal([]byte(containerSettingsJson), &ret)
	if err != nil {
		log.Fatalf("Failed to create container from 'test.container.settings' label. Check JSON format, for detailed format look at https://docs.docker.com/reference/api/docker_remote_api_v1.16/#create-a-container (%s)", containerSettingsJson)
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

func bindEnv(containerSettings *docker.Container, image docker.APIImages) {

	file := EnvFile(image)
	log.Infof("About to bind host Env file %s", file)
	if file != "" {
		if (containerSettings.Config == nil) {
			containerSettings.Config = &docker.Config{}
		}
		common.ReadFile(file, func(line string) {
			appendConfigEnv(containerSettings, line)
		})
		log.Infof("Config env: %+v", containerSettings.Config.Env)
	}
}

func appendConfigEnv(containerSettings *docker.Container, envVar string) {

	containerSettings.Config.Env = append(containerSettings.Config.Env, envVar)
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