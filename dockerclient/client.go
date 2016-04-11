package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)


// A Client is the interface through which mr-burns interacts with the Docker API.
type BurnsDockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]Container, error)
	ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error)
	StartContainer(Container) error
	RemoveContainer(string, bool) error
	WaitContainer(container string) (int, error)
}

// NewClient returns a new Client instance which can be used to interact with
// the Docker API.
func NewClientWithTLS(dockerHost string, pullImages bool, cert, key, ca string) BurnsDockerClient {
	docker, err := docker.NewTLSClient(dockerHost, cert, key, ca)

	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}

	return DockerClient{api: docker, pullImages: pullImages}
}

func NewClient(dockerHost string, pullImages bool) BurnsDockerClient {
	docker, err := docker.NewClient(dockerHost)

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

func (client DockerClient) RemoveContainer(container string, force bool) error {

	log.Infof("Removing container %s", container)
	return client.api.RemoveContainer(docker.RemoveContainerOptions{container, true, force})
}

func (client DockerClient) WaitContainer(container string) (int, error) {

	return client.api.WaitContainer(container)
}

func (client DockerClient) ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error) {

	log.Infof("Retrieving images according to opts: %+v", opts)

	ret, err := client.api.ListImages(opts)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
