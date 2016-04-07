package dockerclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/fsouza/go-dockerclient"
)

func TestID(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{ID: "foo"},
	}

	assert.Equal(t, "foo", c.ID())
}

func TestName(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{Name: "foo"},
	}

	assert.Equal(t, "foo", c.Name())
}

func TestImageID(t *testing.T) {
	c := Container{
		imageInfo: &docker.Image{
			ID: "foo",
		},
	}

	assert.Equal(t, "foo", c.ImageID())
}

func TestImageName_Tagged(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			Image: "foo:latest",
		},
	}

	assert.Equal(t, "foo:latest", c.ImageName())
}

func TestImageName_Untagged(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			Image: "foo",
		},
	}

	assert.Equal(t, "foo:latest", c.ImageName())
}

func TestLinks(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			HostConfig: &docker.HostConfig{
				Links: []string{"foo:foo", "bar:bar"},
			},
		},
	}

	links := c.Links()

	assert.Equal(t, []string{"foo", "bar"}, links)
}

func TestIsTest_True(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			Config: &docker.Config{
				Labels: map[string]string{"test": "true"},
			},
		},
	}

	assert.True(t, c.IsTest())
}

func TestIsTest_WrongLabelValue(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			Config: &docker.Config{
				Labels: map[string]string{"test": "false"},
			},
		},
	}

	assert.False(t, c.IsTest())
}

func TestIsTest_NoLabel(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			Config: &docker.Config{
				Labels: map[string]string{},
			},
		},
	}

	assert.False(t, c.IsTest())
}

func TestRunInterval_Present(t *testing.T) {
	c := Container{
		containerInfo: &docker.Container{
			Config: &docker.Config{
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
		containerInfo: &docker.Container{
			Config: &docker.Config{
				Labels: map[string]string{},
			},
		},
	}

	assert.Equal(t, "", c.RunInterval())
}
