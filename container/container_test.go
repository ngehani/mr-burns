package container

import (
	"testing"

	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{Id: "foo"},
	}

	assert.Equal(t, "foo", c.ID())
}

func TestName(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{Name: "foo"},
	}

	assert.Equal(t, "foo", c.Name())
}

func TestImageID(t *testing.T) {
	c := Container{
		imageInfo: &dockerclient.ImageInfo{
			Id: "foo",
		},
	}

	assert.Equal(t, "foo", c.ImageID())
}

func TestImageName_Tagged(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Image: "foo:latest",
			},
		},
	}

	assert.Equal(t, "foo:latest", c.ImageName())
}

func TestImageName_Untagged(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Image: "foo",
			},
		},
	}

	assert.Equal(t, "foo:latest", c.ImageName())
}

func TestLinks(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			HostConfig: &dockerclient.HostConfig{
				Links: []string{"foo:foo", "bar:bar"},
			},
		},
	}

	links := c.Links()

	assert.Equal(t, []string{"foo", "bar"}, links)
}

func TestIsTest_True(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Labels: map[string]string{"test": "true"},
			},
		},
	}

	assert.True(t, c.IsTest())
}

func TestIsTest_WrongLabelValue(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Labels: map[string]string{"test": "false"},
			},
		},
	}

	assert.False(t, c.IsTest())
}

func TestIsTest_NoLabel(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Labels: map[string]string{},
			},
		},
	}

	assert.False(t, c.IsTest())
}

func TestRunInterval_Present(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Labels: map[string]string{
					"test.run.interval": "300000",
				},
			},
		},
	}

	assert.Equal(t, "300000", c.RunInterval())
}

func TestRunInterval_NoLabel(t *testing.T) {
	c := Container{
		containerInfo: &dockerclient.ContainerInfo{
			Config: &dockerclient.ContainerConfig{
				Labels: map[string]string{},
			},
		},
	}

	assert.Equal(t, "", c.RunInterval())
}
