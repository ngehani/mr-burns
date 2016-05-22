package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
	"strings"
	"fmt"
)

const (
	LABEL_TEST = "test"
	LABEL_INTERVAL = "test.run.interval"
	LABEL_TEST_RESULTS_DIR = "test.results.dir"
	LABEL_TEST_RESULTS_FILE = "test.results.file"
	LABEL_TEST_CONTAINER_SETTINGS = "test.container.settings"
)

// Container represents a running Docker container.
type Container struct {
	Data *docker.Container
}

// NewContainer returns a new Container instance instantiated with the
// specified ContainerInfo and ImageInfo structs.
func NewContainer(data *docker.Container) Container {

	return Container{
		Data: data,
	}
}

// IsTest returns a boolean flag indicating whether or not the current
// container is a test container.
func (c Container) IsTest() bool {

	val, ok := c.Data.Config.Labels[LABEL_TEST]

	return ok && val == "true"
}

func (c Container) RunInterval() string {

	if val, ok := c.Data.Config.Labels[LABEL_INTERVAL]; ok {
		return val
	}

	return ""
}
// Any links in the HostConfig need to be re-written before they can be
// re-submitted to the Docker create API.
func (c Container) GetHostConfig() *docker.HostConfig {

	ret := c.Data.HostConfig
	for i, link := range ret.Links {
		name := link[0:strings.Index(link, ":")]
		alias := link[strings.LastIndex(link, "/"):]
		ret.Links[i] = fmt.Sprintf("%s:%s", name, alias)
	}

	return ret
}
