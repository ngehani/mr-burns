package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-adm/go-dockerclient"
)


// A Client is the interface through which mr-burns interacts with the Docker API.
type DockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]Container, error)
	ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error)
	StartContainer(Container) error
	RemoveContainer(string, bool) error
	WaitContainer(container string) (int, error)
}

type DockerClientContainer struct {
	api        Client
	pullImages bool
}

// NewClient returns a new Client instance which can be used to interact with
// the Docker API.
func NewClientWithTLS(dockerHost string, pullImages bool, cert, key, ca string) DockerClient {
	docker, err := docker.NewTLSClient(dockerHost, cert, key, ca)

	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}

	return DockerClientContainer{api: docker, pullImages: pullImages}
}

func NewClient(dockerHost string, pullImages bool) DockerClient {
	docker, err := docker.NewClient(dockerHost)
	log.Infof("Docker client: %+v", docker)

	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}

	return DockerClientContainer{api: docker, pullImages: pullImages}
}

func (client DockerClientContainer) ListContainers(opts docker.ListContainersOptions) ([]Container, error) {

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

func (client DockerClientContainer) StartContainer(c Container) error {

	log.Infof("Starting %s", c.Name())

	container, err := client.api.CreateContainer(docker.CreateContainerOptions{c.Name(), c.Config(), c.hostConfig()})
	if err != nil {
		return err
	}

	log.Debugf("Starting container %s (%+v)", c.Name(), container)

	return client.api.StartContainer(container.ID, c.hostConfig())
}

func (client DockerClientContainer) RemoveContainer(container string, force bool) error {

	log.Infof("Removing container %s", container)
	return client.api.RemoveContainer(docker.RemoveContainerOptions{container, true, force})
}

func (client DockerClientContainer) WaitContainer(container string) (int, error) {

	return client.api.WaitContainer(container)
}

func (client DockerClientContainer) ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error) {

	log.Infof("Retrieving images according to opts: %+v", opts)

	ret, err := client.api.ListImages(opts)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
