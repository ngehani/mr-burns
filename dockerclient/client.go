package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

// A Client is the interface through which mr-burns interacts with the Docker API.
type DockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]Container, error)
	ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error)
	StartContainer(Container) error
	RemoveContainer(string, bool) error
	WaitContainer(container string) (int, error)
}

type DockerClientWrapper struct {
	client Client
}

func NewClient(dockerHost string) DockerClient {

	docker, err := docker.NewClient(dockerHost)
	log.Infof("Docker client: %+v", docker)
	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}

	return DockerClientWrapper{client: docker}
}

func (wrapper DockerClientWrapper) ListContainers(opts docker.ListContainersOptions) ([]Container, error) {

	ret := []Container{}
	log.Infof("Retrieving containers according to opts: %+v", opts)
	cs, err := wrapper.client.ListContainers(opts)
	if err != nil {
		return nil, err
	}

	for _, container := range cs {
		containerInfo, err := wrapper.client.InspectContainer(container.ID)
		if err != nil {
			return nil, err
		}

		imageInfo, err := wrapper.client.InspectImage(containerInfo.Image)
		if err != nil {
			return nil, err
		}

		ret = append(ret, Container{containerInfo: containerInfo, imageInfo: imageInfo})
	}

	return ret, nil
}

func (wrapper DockerClientWrapper) StartContainer(c Container) error {

	log.Infof("Creating container %s", c.Name())
	container, err := wrapper.client.CreateContainer(docker.CreateContainerOptions{c.Name(), c.Config(), c.hostConfig()})
	if err != nil {
		return err
	}
	log.Infof("Starting container %s (%+v)", c.Name(), container)

	return wrapper.client.StartContainer(container.ID, c.hostConfig())
}

func (wrapper DockerClientWrapper) RemoveContainer(container string, force bool) error {

	log.Infof("Removing container %s", container)
	return wrapper.client.RemoveContainer(docker.RemoveContainerOptions{container, true, force})
}

func (wrapper DockerClientWrapper) WaitContainer(container string) (int, error) {

	return wrapper.client.WaitContainer(container)
}

func (wrapper DockerClientWrapper) ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error) {

	log.Infof("Retrieving images according to opts: %+v", opts)
	ret, err := wrapper.client.ListImages(opts)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
