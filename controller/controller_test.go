package controller

import (
	"testing"
	"github.com/gaia-adm/mr-burns/dockerclient"
	"github.com/gaia-adm/go-dockerclient"
	"github.com/stretchr/testify/mock"
)

func TestStart(t *testing.T) {

	mockClient := dockerclient.NewMockClient()
	labels := map[string]string{dockerclient.TestResultsLabel: "/tmp/junit-results"}
	img := docker.APIImages{Labels:labels}
	images := []docker.APIImages{img}
	mockClient.On("ListImages", mock.Anything).Return(images)
	Start(mockClient)
}
