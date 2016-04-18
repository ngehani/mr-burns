package controller

import (
	"testing"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/mock"
)

func TestStart(t *testing.T) {

	mockClient := dockerclient.NewMockClient()
	images := addImages(mockClient)
	addRemoveContainer(mockClient)
	addCreateContainer(mockClient, images)
	addStartContainer(mockClient)
	addWaitContainer(mockClient)
	Start(mockClient.CreateMockClientWrapper())
}

func addImages(mockClient *dockerclient.MockClient) []docker.APIImages {

	labels := map[string]string{dockerclient.LabelTestResultPath: "/tmp/junit-results"}
	images := []docker.APIImages{docker.APIImages{Labels:labels}}
	mockClient.On("ListImages", mock.Anything).Return(images, nil)

	return images
}

func addRemoveContainer(mockClient *dockerclient.MockClient) {

	mockClient.On("RemoveContainer", mock.Anything, mock.Anything).Return(nil)
}

func addCreateContainer(mockClient *dockerclient.MockClient, images []docker.APIImages) {

	var options docker.CreateContainerOptions
	mockClient.On("CreateContainer", mock.Anything).
	Run(func(args mock.Arguments) {
		options = args.Get(0).(docker.CreateContainerOptions)
	}).
	Return(&docker.Container{
		Name: options.Name,
		Config: options.Config,
		HostConfig: options.HostConfig},
		nil)
}

func addStartContainer(mockClient *dockerclient.MockClient) {

	mockClient.On("StartContainer", mock.Anything, mock.Anything).Return(nil)
}

func addWaitContainer(mockClient *dockerclient.MockClient) {

	mockClient.On("WaitContainer", mock.Anything).Return(-1, nil)
}