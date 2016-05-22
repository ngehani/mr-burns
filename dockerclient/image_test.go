package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunInterval_Present(t *testing.T) {

	const INTERVAL = "300000"
	image := docker.APIImages{ID: "my-test-image", Labels: map[string]string{LABEL_TEST_RUN_INTERVAL: INTERVAL}}

	assert.Equal(t, INTERVAL, RunInterval(image))
}

func TestRunInterval_NoLabel(t *testing.T) {

	image := docker.APIImages{ID: "my-test-image"}
	assert.Equal(t, "", RunInterval(image))
}