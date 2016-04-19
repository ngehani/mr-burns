package controller

import (
	"testing"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"sync"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestStartController(t *testing.T) {

	const IMAGE_ID string = "aa789bb"
	client := dockerclient.NewMockClient()
	manager := NewDockerManager(client.CreateMockClientWrapper())
	waitToPublish := sync.WaitGroup{}
	waitToPublish.Add(1)
	sleepTimes := 1
	var givenResultsFilePath string
	controller := Controller{
		taskIdToTask: map[string]Task{},
		docker: manager,
		sleep: func() {
			waitToPublish.Wait()
		},
		publish: func(resultsFilePath string) error {
			defer waitToPublish.Done()
			givenResultsFilePath = resultsFilePath
			return nil
		},
		stop: func() bool {
			if sleepTimes == 1 {
				sleepTimes++
				return false
			}
			return true
		},
	}
	mockListImages(client, IMAGE_ID)
	mockRemoveContainer(client)
	mockCreateContainer(client)
	mockStartContainer(client)
	mockWaitContainer(client)
	controller.Start()

	assert.True(t, strings.Contains(givenResultsFilePath, IMAGE_ID),
		"Given results file path does not contains image ID. Path:", givenResultsFilePath,
		"Image ID:", IMAGE_ID)
}
