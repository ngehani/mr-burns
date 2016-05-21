package controller

import (
	"testing"
	"github.com/fsouza/go-dockerclient"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/stretchr/testify/assert"
	"strings"
)

const (
	IMAGE_ID string = "test-image-id"
	RESULTS_DIR = "/src/results"
)

func TestBuildContainer(t *testing.T) {

	container := BuildContainer(getImage(getContainerSettingsMockJson()), "test-container-name")
	assert.Equal(t, IMAGE_ID, container.Config.Image)
	assert.Equal(t, "Effi", container.Config.User)
	assert.True(t, arrContainsSubString(container.HostConfig.Binds, RESULTS_DIR))
}

func TestBuildContainerEmptyContainerSettings(t *testing.T) {

	container := BuildContainer(getImage(""), "test-container-name")
	assert.True(t, arrContainsSubString(container.HostConfig.Binds, RESULTS_DIR))
}

func getImage(containerSettingsJson string) docker.APIImages {

	return docker.APIImages{ID: IMAGE_ID, Labels: map[string]string{
		dockerclient.LABEL_TEST_CONTAINER_SETTINGS: containerSettingsJson,
		dockerclient.LABEL_TEST_RESULTS_DIR: RESULTS_DIR}}
}

func getContainerSettingsMockJson() string {

	return `{
             "Id": "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2",
             "Created": "2013-05-07T14:51:42.087658+02:00",
             "Path": "date",
             "Args": [],
             "Config": {
                     "Hostname": "4fa6e0f0c678",
                     "User": "Effi",
                     "Memory": 17179869184,
                     "MemorySwap": -1,
                     "AttachStdin": false,
                     "AttachStdout": true,
                     "AttachStderr": true,
                     "PortSpecs": null,
                     "Tty": false,
                     "OpenStdin": false,
                     "StdinOnce": false,
                     "Env": null,
                     "Cmd": [
                             "date"
                     ],
                     "Image": "base",
                     "Volumes": {},
                     "VolumesFrom": ""
             },
             "State": {
                     "Running": false,
                     "Pid": 0,
                     "ExitCode": 0,
                     "StartedAt": "2013-05-07T14:51:42.087658+02:00",
                     "Ghost": false
             },
             "Image": "b750fe79269d2ec9a3c593ef05b4332b1d1a02a62b4accb2c21d589ff2f5f2dc",
             "NetworkSettings": {
                     "IpAddress": "",
                     "IpPrefixLen": 0,
                     "Gateway": "",
                     "Bridge": "",
                     "PortMapping": null
             },
             "SysInitPath": "/home/kitty/go/src/github.com/dotcloud/docker/bin/docker",
             "ResolvConfPath": "/etc/resolv.conf",
             "Volumes": {},
             "HostConfig": {
               "Binds": ["/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt"],
               "ContainerIDFile": "",
               "LxcConf": [],
               "Privileged": false,
               "PortBindings": {
                 "80/tcp": [
                   {
                     "HostIp": "0.0.0.0",
                     "HostPort": "49153"
                   }
                 ]
               },
               "Links": null,
               "PublishAllPorts": false
             }
	}`
}

func arrContainsSubString(arr []string, substring string) bool {

	ret := false
	for _, curr := range arr {
		if strings.Contains(curr, substring) {
			ret = true
			break
		}
	}

	return ret
}