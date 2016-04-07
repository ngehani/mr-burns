package container

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)


// A Client is the interface through which mr-burns interacts with the Docker API.
type burnsDockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]Container, error)
	StartContainer(Container) error
	RemoveContainer(Container, bool) error
}

// NewClient returns a new Client instance which can be used to interact with
// the Docker API.
func NewClient(dockerHost string, pullImages bool, cert, key, ca string) burnsDockerClient {
	docker, err := docker.NewTLSClient(dockerHost, cert, key, ca)

	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}

	return DockerClient{api: docker, pullImages: pullImages}
}

type DockerClient struct {
	api        Client
	pullImages bool
}

func (client DockerClient) ListContainers(opts docker.ListContainersOptions) ([]Container, error) {

	ret := []Container{}
	log.Infof("Retrieving containers according to opts: %+v", opts)

	cs, err := client.api.ListContainers(opts)
	if err != nil {
		return nil, err
	}

	for _, container := range cs {
		containerInfo, err := client.api.InspectContainer(container.ID)
		if err != nil {
			return nil, err
		}

		imageInfo, err := client.api.InspectImage(containerInfo.Image)
		if err != nil {
			return nil, err
		}

		ret = append(ret, Container{containerInfo: containerInfo, imageInfo: imageInfo})
	}

	return ret, nil
}

func (client DockerClient) StartContainer(c Container) error {

	log.Infof("Starting %s", c.Name())

	container, err := client.api.CreateContainer(docker.CreateContainerOptions{c.Name(), c.Config(), c.hostConfig()})
	if err != nil {
		return err
	}

	log.Debugf("Starting container %s (%+v)", c.Name(), container)

	return client.api.StartContainer(container.ID, c.hostConfig())
}

func (client DockerClient) RemoveContainer(c Container, force bool) error {

	log.Infof("Removing container %s", c.ID())
	return client.api.RemoveContainer(docker.RemoveContainerOptions{c.ID(), true, force})
}
