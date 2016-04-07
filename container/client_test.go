package container

import (
	"testing"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/gaia-adm/mr-burns/container/mockclient"
	"errors"
	"github.com/stretchr/testify/mock"
)

func TestListContainers_Success(t *testing.T) {
	ci := createDummyContainerInfo()
	api := mockclient.NewMockClient()
	ii := &docker.Image{}
	lco := docker.ListContainersOptions{All:false, Size:false}
	api.On("ListContainers", lco).Return([]docker.APIContainers{{ID: "foo", Names:[]string{"bar"}}}, nil)
	api.On("InspectContainer", "foo").Return(ci, nil)
	api.On("InspectImage", "abc123").Return(ii, nil)

	client := DockerClient{api: api}
	cs, err := client.ListContainers(lco)

	assert.NoError(t, err)
	assert.Len(t, cs, 1)
	assert.Equal(t, ci, cs[0].containerInfo)
	assert.Equal(t, ii, cs[0].imageInfo)
}

func TestListContainers_Filter(t *testing.T) {
	ci := createDummyContainerInfo()
	api := mockclient.NewMockClient()
	ii := &docker.Image{}
	lco := docker.ListContainersOptions{All:false, Size:false, Filters: map[string][]string{"label": {"test="}}}
	api.On("ListContainers", lco).Return([]docker.APIContainers{{ID: "foo", Names:[]string{"bar"}, Labels: map[string]string{"label": "test="}}}, nil)
	api.On("InspectContainer", "foo").Return(ci, nil)
	api.On("InspectImage", "abc123").Return(ii, nil)

	client := DockerClient{api: api}
	cs, err := client.ListContainers(lco)

	assert.NoError(t, err)
	assert.Len(t, cs, 1)
	assert.Equal(t, ci, cs[0].containerInfo)
	assert.Equal(t, ii, cs[0].imageInfo)
}

func TestListContainers_ListError(t *testing.T) {
	api := mockclient.NewMockClient()
	lco := docker.ListContainersOptions{All:false, Size:false}
	api.On("ListContainers", lco).Return([]docker.APIContainers{}, errors.New("oops"))

	client := DockerClient{api: api}
	_, err := client.ListContainers(lco)

	assert.Error(t, err)
	assert.EqualError(t, err, "oops")
	api.AssertExpectations(t)
}

func TestListContainers_InspectContainerError(t *testing.T) {
	api := mockclient.NewMockClient()
	lco := docker.ListContainersOptions{All:false, Size:false}
	api.On("ListContainers", lco).Return([]docker.APIContainers{{ID: "foo", Names:[]string{"bar"}}}, nil)
	api.On("InspectContainer", "foo").Return(&docker.Container{}, errors.New("uh-oh"))

	client := DockerClient{api: api}
	_, err := client.ListContainers(lco)

	assert.Error(t, err)
	assert.EqualError(t, err, "uh-oh")
	api.AssertExpectations(t)
}

func TestListContainers_InspectImageError(t *testing.T) {
	ci := &docker.Container{Image: "abc123", Config: &docker.Config{Image: "img"}}
	ii := &docker.Image{}
	lco := docker.ListContainersOptions{All:false, Size:false}
	api := mockclient.NewMockClient()
	api.On("ListContainers", lco).Return([]docker.APIContainers{{ID: "foo", Names:[]string{"bar"}}}, nil)
	api.On("InspectContainer", "foo").Return(ci, nil)
	api.On("InspectImage", "abc123").Return(ii, errors.New("whoops"))

	client := DockerClient{api: api}
	_, err := client.ListContainers(lco)

	assert.Error(t, err)
	assert.EqualError(t, err, "whoops")
	api.AssertExpectations(t)
}

func createDummyContainerInfo() *docker.Container {

	return &docker.Container{Image: "abc123", Config: &docker.Config{Image: "img"}}
}

func TestStartContainer_Success(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			ID: "def789",
			Name:       "foo",
			Config:     &docker.Config{},
			HostConfig: &docker.HostConfig{},
		},
		imageInfo: &docker.Image{
			Config: &docker.Config{},
		},
	}

	api := mockclient.NewMockClient()
	api.On("CreateContainer",
		mock.Anything).Return(c.containerInfo, nil)
	api.On("StartContainer", "def789", mock.Anything).Return(nil)

	client := DockerClient{api: api}
	err := client.StartContainer(c)

	assert.NoError(t, err)
	api.AssertExpectations(t)
}

func TestStartContainer_CreateContainerError(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			ID: "def789",
			Name:       "foo",
			Config:     &docker.Config{},
			HostConfig: &docker.HostConfig{},
		},
	}

	api := mockclient.NewMockClient()
	api.On("CreateContainer", mock.Anything).Return(c.containerInfo, errors.New("oops"))

	client := DockerClient{api: api}
	err := client.StartContainer(c)

	assert.Error(t, err)
	assert.EqualError(t, err, "oops")
	api.AssertExpectations(t)
}

func TestStartContainer_StartContainerError(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			ID: "def789",
			Name:       "foo",
			Config:     &docker.Config{},
			HostConfig: &docker.HostConfig{},
		},
	}

	api := mockclient.NewMockClient()
	api.On("CreateContainer", mock.Anything).Return(c.containerInfo, nil)
	api.On("StartContainer", "def789", mock.Anything).Return(errors.New("whoops"))

	client := DockerClient{api: api}
	err := client.StartContainer(c)

	assert.Error(t, err)
	assert.EqualError(t, err, "whoops")
	api.AssertExpectations(t)
}


