package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsTest_True(t *testing.T) {

	c := NewContainer(&docker.Container{
		Config: &docker.Config{
			Labels: map[string]string{LABEL_TEST: "true"},
		},
	})

	assert.True(t, c.IsTest())
}

func TestIsTest_WrongLabelValue(t *testing.T) {

	c := NewContainer(&docker.Container{
		Config: &docker.Config{
			Labels: map[string]string{LABEL_TEST: "false"},
		},
	})

	assert.False(t, c.IsTest())
}

func TestIsTest_NoLabel(t *testing.T) {

	c := NewContainer(&docker.Container{
		Config: &docker.Config{
			Labels: map[string]string{},
		},
	})

	assert.False(t, c.IsTest())
}

func TestRunInterval_Present(t *testing.T) {

	c := NewContainer(&docker.Container{
		Config: &docker.Config{
			Labels: map[string]string{
				LABEL_INTERVAL: "300000",
			},
		},
	})

	assert.Equal(t, "300000", c.RunInterval())
}

func TestRunInterval_NoLabel(t *testing.T) {

	c := NewContainer(&docker.Container{
		Config: &docker.Config{
			Labels: map[string]string{},
		},
	})

	assert.Equal(t, "", c.RunInterval())
}