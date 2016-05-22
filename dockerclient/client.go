package dockerclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
	"bytes"
	"bufio"
)

type DockerClient struct {
	client Client
}

func NewClient(dockerHost string) DockerClient {

	docker, err := docker.NewClient(dockerHost)
	log.Infof("Docker client: %+v", docker)
	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}

	return DockerClient{client: docker}
}

func (client DockerClient) StartContainer(c Container) error {

	log.Infof("Creating container %s", c.Data.Name)
	container, err := client.client.CreateContainer(docker.CreateContainerOptions{c.Data.Name, c.Data.Config, c.HostConfig()})
	if err != nil {
		return err
	}
	log.Infof("Starting container %s (%+v)", c.Data.Name, container)

	return client.client.StartContainer(container.ID, c.HostConfig())
}

func (client DockerClient) RemoveContainer(container string, force bool) error {

	log.Infof("Removing container %s", container)
	return client.client.RemoveContainer(docker.RemoveContainerOptions{container, true, force})
}

func (client DockerClient) WaitContainer(container string) (int, error) {

	return client.client.WaitContainer(container)
}

func (client DockerClient) ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error) {

	log.Infof("Retrieving docker images according to opts: %+v", opts)
	ret, err := client.client.ListImages(opts)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (client DockerClient) Logs(container string) (string, error) {

	reader, writer := io.Pipe()
	err := client.client.Logs(docker.LogsOptions{
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