package dockerclient

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/fsouza/go-dockerclient"
)

func TestStartContainer_Success(t *testing.T) {

	c := NewContainer(&docker.Container{
		ID: "def789",
		Name:       "foo",
		Config:     &docker.Config{},
		HostConfig: &docker.HostConfig{},
	})
	api := NewMockClient()
	api.On("CreateContainer",
		mock.Anything).Return(c.Data, nil)
	api.On("StartContainer", "def789", mock.Anything).Return(nil)

	client := api.CreateMockClientWrapper()
	err := client.StartContainer(c)

	assert.NoError(t, err)
	api.AssertExpectations(t)
}

func TestStartContainer_CreateContainerError(t *testing.T) {

	c := NewContainer(&docker.Container{
		ID: "def789",
		Name:       "foo",
		Config:     &docker.Config{},
		HostConfig: &docker.HostConfig{},
	})
	api := NewMockClient()
	api.On("CreateContainer", mock.Anything).Return(c.Data, errors.New("oops"))

	client := api.CreateMockClientWrapper()
	err := client.StartContainer(c)

	assert.Error(t, err)
	assert.EqualError(t, err, "oops")
	api.AssertExpectations(t)
}

func TestStartContainer_StartContainerError(t *testing.T) {

	c := NewContainer(&docker.Container{
		ID: "def789",
		Name:       "foo",
		Config:     &docker.Config{},
		HostConfig: &docker.HostConfig{},
	})

	api := NewMockClient()
	api.On("CreateContainer", mock.Anything).Return(c.Data, nil)
	api.On("StartContainer", "def789", mock.Anything).Return(errors.New("whoops"))

	client := api.CreateMockClientWrapper()
	err := client.StartContainer(c)

	assert.Error(t, err)
	assert.EqualError(t, err, "whoops")
	api.AssertExpectations(t)
}

func TestRemoveContainer(t *testing.T) {

	api := NewMockClient()
	api.On("RemoveContainer", docker.RemoveContainerOptions{"foo", true, true}).Return(nil)

	client := api.CreateMockClientWrapper()
	err := client.RemoveContainer("foo", true)

	assert.NoError(t, err)
	api.AssertExpectations(t)
}