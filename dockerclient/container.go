package dockerclient

import (
	"fmt"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

const (
	testLabel = "test"
	intervalLabel = "test.run.interval"
)

// NewContainer returns a new Container instance instantiated with the
// specified ContainerInfo and ImageInfo structs.
func NewContainer(containerInfo *docker.Container, imageInfo *docker.Image) *Container {
	return &Container{
		containerInfo: containerInfo,
		imageInfo:     imageInfo,
	}
}

// Container represents a running Docker container.
type Container struct {
	containerInfo *docker.Container
	imageInfo     *docker.Image
}

// ID returns the Docker container ID.
func (c Container) ID() string {
	return c.containerInfo.ID
}

// Name returns the Docker container name.
func (c Container) Name() string {
	return c.containerInfo.Name
}

// ImageID returns the ID of the Docker image that was used to start the container.
func (c Container) ImageID() string {
	return c.imageInfo.ID
}

// ImageName returns the name of the Docker image that was used to start the
// container. If the original image was specified without a particular tag, the
// "latest" tag is assumed.
func (c Container) ImageName() string {
	imageName := c.containerInfo.Image
	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	return imageName
}

// Links returns a list containing the names of all the containers to which
// this container is linked.
func (c Container) Links() []string {
	var links []string

	if (c.containerInfo != nil) && (c.containerInfo.HostConfig != nil) {
		for _, link := range c.containerInfo.HostConfig.Links {
			name := strings.Split(link, ":")[0]
			links = append(links, name)
		}
	}

	return links
}

// IsTest returns a boolean flag indicating whether or not the current
// container is a test container.
func (c Container) IsTest() bool {
	val, ok := c.containerInfo.Config.Labels[testLabel]
	return ok && val == "true"
}

func (c Container) RunInterval() string {
	if val, ok := c.containerInfo.Config.Labels[intervalLabel]; ok {
		return val
	}

	return ""
}

func (c Container) Config() *docker.Config {
	config := c.containerInfo.Config

	return config
}

// Any links in the HostConfig need to be re-written before they can be
// re-submitted to the Docker create API.
func (c Container) hostConfig() *docker.HostConfig {
	hostConfig := c.containerInfo.HostConfig

	for i, link := range hostConfig.Links {
		name := link[0:strings.Index(link, ":")]
		alias := link[strings.LastIndex(link, "/"):]

		hostConfig.Links[i] = fmt.Sprintf("%s:%s", name, alias)
	}

	return hostConfig
}
