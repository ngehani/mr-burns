package controller

import (
	"testing"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const (
	MOCK_TEST_RESULTS_FILE_NAME string = "mock-test-results.xml"
	MOCK_TEST_DESC string = "Simpsons Integration Tests"
)

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

func TestGetImagesFilterDangling(t *testing.T) {

	const IMAGE_ID string = "aa789bb"
	client := dockerclient.NewMockClient()
	manager := NewDockerManager(client.CreateMockClientWrapper())
	mockListImagesRepoTags(client, IMAGE_ID, "<none>:<none>")
	images, err := manager.GetImages()

	assert.NoError(t, err)
	assert.Len(t, images, 0)
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

	return mockListImagesRepoTags(mockClient, imageId, "gaiaadm/mr-burns-builder:latest")
}

func mockListImagesRepoTags(mockClient *dockerclient.MockClient, imageId string, repoTags string) docker.APIImages {

	labels := map[string]string{
		dockerclient.LABEL_TEST_RESULTS_DIR: "/tmp/test-results",
		dockerclient.LABEL_TEST_RESULTS_FILE: MOCK_TEST_RESULTS_FILE_NAME,
		dockerclient.LABEL_INTERVAL: "600000",
		dockerclient.LABEL_DESC: MOCK_TEST_DESC}
	ret := docker.APIImages{ID: imageId, Labels: labels, RepoTags: []string{repoTags}}
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

func mockLogs(mockClient *dockerclient.MockClient) string {

	logs := "mkdir: can't create directory '.cover': File exists\n[INFO] Downloading dependencies. Please wait...\n[INFO] Fetching updates for github.com/davecgh/go-spew."
	mockClient.On("Logs", mock.Anything).Return(nil)

	return logs
}