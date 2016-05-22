package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
	"strings"
	"fmt"
)

// Container represents a running Docker container.
type Container struct {
	Data *docker.Container
}

// NewContainer returns a new Container instance instantiated with the
// specified ContainerInfo and ImageInfo structs.
func NewContainer(data *docker.Container) Container {

	return Container{Data: data, }
}

// Any links in the HostConfig need to be re-written before they can be
// re-submitted to the Docker create API.
func (c Container) HostConfig() *docker.HostConfig {

	ret := c.Data.HostConfig
	for i, link := range ret.Links {
		name := link[0:strings.Index(link, ":")]
		alias := link[strings.LastIndex(link, "/"):]
		ret.Links[i] = fmt.Sprintf("%s:%s", name, alias)
	}

	return ret
}
