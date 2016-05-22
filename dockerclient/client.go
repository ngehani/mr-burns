package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
	"bytes"
	"bufio"
)

// A Client is the interface through which mr-burns interacts with the Docker API.
type DockerClient interface {
	ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error)
	StartContainer(Container) error
	RemoveContainer(string, bool) error
	WaitContainer(string) (int, error)
	Logs(container string) (string, error)
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

func (wrapper DockerClientWrapper) StartContainer(c Container) error {

	log.Infof("Creating container %s", c.Data.Name)
	container, err := wrapper.client.CreateContainer(docker.CreateContainerOptions{c.Data.Name, c.Data.Config, c.GetHostConfig()})
	if err != nil {
		return err
	}
	log.Infof("Starting container %s (%+v)", c.Data.Name, container)

	return wrapper.client.StartContainer(container.ID, c.GetHostConfig())
}

func (wrapper DockerClientWrapper) RemoveContainer(container string, force bool) error {

	log.Infof("Removing container %s", container)
	return wrapper.client.RemoveContainer(docker.RemoveContainerOptions{container, true, force})
}

func (wrapper DockerClientWrapper) WaitContainer(container string) (int, error) {

	return wrapper.client.WaitContainer(container)
}

func (wrapper DockerClientWrapper) ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error) {

	log.Infof("Retrieving docker images according to opts: %+v", opts)
	ret, err := wrapper.client.ListImages(opts)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (wrapper DockerClientWrapper) Logs(container string) (string, error) {

	reader, writer := io.Pipe()
	err := wrapper.client.Logs(docker.LogsOptions{
		Container: container,
		OutputStream: writer,
		ErrorStream:  writer,
		Stdout: true,
		Stderr: true,
		Timestamps:   true,
	})
	if err != nil {
		return "", err
	}

	logs := make(chan string)
	errScanner := make(chan error)
	// read stdout and stderr logs
	go func(reader io.Reader) {

		var buffer bytes.Buffer
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			buffer.WriteString(scanner.Text())
		}
		logs <- buffer.String()
		errScanner <- scanner.Err()
	}(reader)
	ret := <-logs
	errScannerReceiver := <-errScanner

	return ret, errScannerReceiver
}