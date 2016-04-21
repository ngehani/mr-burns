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
	const REPO_TAGS string = "check:me"
	const RESULTS string = "Mock results - Tests run: 2, Failures: 0, Errors: 0, Skipped: 0, Time elapsed: 5.046 sec"
	client := dockerclient.NewMockClient()
	manager := NewDockerManager(client.CreateMockClientWrapper())
	waitToPublish := sync.WaitGroup{}
	waitToPublish.Add(1)
	stopTimes := 0
	var publishData string
	controller := Controller{
		taskIdToTask: map[string]Task{},
		docker: manager,
		sleep: func() {
			waitToPublish.Wait()
		},
		getResults: func(string) string {
			return RESULTS
		},
		publish: func(data string) error {
			defer waitToPublish.Done()
			publishData = data
			return nil
		},
		stop: func() bool {
			stopTimes++
			if stopTimes == 1 {
				return false
			}
			return true
		},
	}
	mockListImagesRepoTags(client, IMAGE_ID, REPO_TAGS)
	mockRemoveContainer(client)
	mockCreateContainer(client)
	mockStartContainer(client)
	mockWaitContainer(client)
	mockLogs(client)
	controller.Start()

	assert.Equal(t, 2, stopTimes)
	assert.True(t, strings.Contains(publishData, REPO_TAGS), publishData)
	assert.True(t, strings.Contains(publishData, RESULTS), publishData)
}
