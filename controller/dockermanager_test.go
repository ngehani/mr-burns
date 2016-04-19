package controller

import (
	"testing"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const MOCK_TEST_RESULTS_FILE_NAME string = "mock-test-results.xml"

func TestGetImages(t *testing.T) {

	const IMAGE_ID string = "aa789bb"
	client := dockerclient.NewMockClient()
	manager := NewDockerManager(client.CreateMockClientWrapper())
	mockImage := mockListImages(client, IMAGE_ID)
	images, err := manager.GetImages()

	assert.NoError(t, err)
	assert.Len(t, images, 1)
	assert.Equal(t, mockImage.ID, images[0].ID)
}

func TestRunTests(t *testing.T) {

	const IMAGE_ID string = "aa789bb"
	client := dockerclient.NewMockClient()
	manager := NewDockerManager(client.CreateMockClientWrapper())
	mockImage := mockListImages(client, IMAGE_ID)
	mockRemoveContainer(client)
	mockCreateContainer(client)
	mockStartContainer(client)
	mockWaitContainer(client)
	_, err := manager.RunTests(mockImage, fmt.Sprintf("%s@mr-burns", IMAGE_ID))

	assert.NoError(t, err)
}

func mockListImages(mockClient *dockerclient.MockClient, imageId string) docker.APIImages {

	labels := map[string]string{
		dockerclient.LabelTestResultPath: "/tmp/test-results",
		dockerclient.LabelTestResultsFile: MOCK_TEST_RESULTS_FILE_NAME,
		dockerclient.LabelInterval: "600000"}
	ret := docker.APIImages{ID: imageId, Labels: labels, RepoTags: []string{"gaiaadm/mr-burns-builder:latest"}}
	mockClient.On("ListImages", mock.Anything).Return([]docker.APIImages{ret}, nil)

	return ret
}

func mockRemoveContainer(mockClient *dockerclient.MockClient) {

	mockClient.On("RemoveContainer", mock.Anything, mock.Anything).Return(nil)
}

func mockCreateContainer(mockClient *dockerclient.MockClient) {

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

func mockStartContainer(mockClient *dockerclient.MockClient) {

	mockClient.On("StartContainer", mock.Anything, mock.Anything).Return(nil)
}

func mockWaitContainer(mockClient *dockerclient.MockClient) {

	mockClient.On("WaitContainer", mock.Anything).Return(0, nil)
}