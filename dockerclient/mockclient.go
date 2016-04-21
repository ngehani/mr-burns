package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/mock"
	"errors"
)

type MockClient struct {
	mock.Mock
}

func NewMockClient() *MockClient {

	return &MockClient{}
}

func (client *MockClient) CreateMockClientWrapper() DockerClientWrapper {

	return DockerClientWrapper{client: client}
}

func (client *MockClient) Endpoint() string {

	args := client.Mock.Called()

	return args.String(0)
}

func (client *MockClient) Ping() error {

	args := client.Mock.Called()

	return args.Error(0)
}

func (client *MockClient) ListVolumes(opts docker.ListVolumesOptions) ([]docker.Volume, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).([]docker.Volume), args.Error(1)
}

func (client *MockClient) CreateVolume(opts docker.CreateVolumeOptions) (*docker.Volume, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).(*docker.Volume), args.Error(1)
}

func (client *MockClient) InspectVolume(name string) (*docker.Volume, error) {

	args := client.Mock.Called(name)

	return args.Get(0).(*docker.Volume), args.Error(1)
}

func (client *MockClient) RemoveVolume(name string) error {

	args := client.Mock.Called(name)

	return args.Error(0)
}

func (client *MockClient) ListNetworks() ([]docker.Network, error) {
	args := client.Mock.Called()
	return args.Get(0).([]docker.Network), args.Error(1)
}

func (client *MockClient) FilteredListNetworks(opts docker.NetworkFilterOpts) ([]docker.Network, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).([]docker.Network), args.Error(1)
}

func (client *MockClient) NetworkInfo(id string) (*docker.Network, error) {

	args := client.Mock.Called(id)

	return args.Get(0).(*docker.Network), args.Error(1)
}

func (client *MockClient) CreateNetwork(opts docker.CreateNetworkOptions) (*docker.Network, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).(*docker.Network), args.Error(1)
}

func (client *MockClient) RemoveNetwork(id string) error {

	args := client.Mock.Called(id)

	return args.Error(0)
}

func (client *MockClient) ConnectNetwork(id string, opts docker.NetworkConnectionOptions) error {
	args := client.Mock.Called(id, opts)
	return args.Error(0)
}

func (client *MockClient) DisconnectNetwork(id string, opts docker.NetworkConnectionOptions) error {

	args := client.Mock.Called(id, opts)

	return args.Error(0)
}

func (client *MockClient) CreateExec(opts docker.CreateExecOptions) (*docker.Exec, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).(*docker.Exec), args.Error(1)
}

func (client *MockClient) StartExec(id string, opts docker.StartExecOptions) error {
	args := client.Mock.Called(id, opts)
	return args.Error(0)
}

func (client *MockClient) StartExecNonBlocking(id string, opts docker.StartExecOptions) (docker.CloseWaiter, error) {

	args := client.Mock.Called(id, opts)
	return args.Get(0).(docker.CloseWaiter), args.Error(1)
}

func (client *MockClient) ResizeExecTTY(id string, height, width int) error {

	args := client.Mock.Called(id, height, width)
	return args.Error(0)
}

func (client *MockClient) InspectExec(id string) (*docker.ExecInspect, error) {

	args := client.Mock.Called(id)
	return args.Get(0).(*docker.ExecInspect), args.Error(1)
}

func (client *MockClient) Version() (*docker.Env, error) {

	args := client.Mock.Called()
	return args.Get(0).(*docker.Env), args.Error(1)
}

func (client *MockClient) Info() (*docker.DockerInfo, error) {

	args := client.Mock.Called()
	return args.Get(0).(*docker.DockerInfo), args.Error(1)
}

func (client *MockClient) ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error) {

	args := client.Mock.Called(opts)
	return args.Get(0).([]docker.APIImages), args.Error(1)
}

func (client *MockClient) ImageHistory(name string) ([]docker.ImageHistory, error) {

	args := client.Mock.Called(name)
	return args.Get(0).([]docker.ImageHistory), args.Error(1)
}

func (client *MockClient) RemoveImage(name string) error {

	args := client.Mock.Called(name)
	return args.Error(0)
}

func (client *MockClient) RemoveImageExtended(name string, opts docker.RemoveImageOptions) error {

	args := client.Mock.Called(name, opts)
	return args.Error(0)
}

func (client *MockClient) InspectImage(name string) (*docker.Image, error) {

	args := client.Mock.Called(name)
	return args.Get(0).(*docker.Image), args.Error(1)
}

func (client *MockClient) PushImage(opts docker.PushImageOptions, auth docker.AuthConfiguration) error {

	args := client.Mock.Called(opts, auth)
	return args.Error(0)
}

func (client *MockClient) PullImage(opts docker.PullImageOptions, auth docker.AuthConfiguration) error {

	args := client.Mock.Called(opts, auth)
	return args.Error(0)
}

func (client *MockClient) LoadImage(opts docker.LoadImageOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) ExportImage(opts docker.ExportImageOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) ExportImages(opts docker.ExportImagesOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) ImportImage(opts docker.ImportImageOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) BuildImage(opts docker.BuildImageOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) TagImage(name string, opts docker.TagImageOptions) error {

	args := client.Mock.Called(name, opts)
	return args.Error(0)
}

func (client *MockClient) SearchImages(term string) ([]docker.APIImageSearch, error) {

	args := client.Mock.Called(term)
	return args.Get(0).([]docker.APIImageSearch), args.Error(1)
}

func (client *MockClient) SearchImagesEx(term string, auth docker.AuthConfiguration) ([]docker.APIImageSearch, error) {

	args := client.Mock.Called(term, auth)
	return args.Get(0).([]docker.APIImageSearch), args.Error(1)
}

func (client *MockClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).([]docker.APIContainers), args.Error(1)
}

func (client *MockClient) UpdateContainer(id string, opts docker.UpdateContainerOptions) error {

	args := client.Mock.Called(id, opts)

	return args.Error(0)
}

func (client *MockClient) RenameContainer(opts docker.RenameContainerOptions) error {

	args := client.Mock.Called(opts)

	return args.Error(0)
}

func (client *MockClient) InspectContainer(id string) (*docker.Container, error) {

	args := client.Mock.Called(id)
	return args.Get(0).(*docker.Container), args.Error(1)
}

func (client *MockClient) ContainerChanges(id string) ([]docker.Change, error) {

	args := client.Mock.Called(id)
	return args.Get(0).([]docker.Change), args.Error(1)
}

func (client *MockClient) CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).(*docker.Container), args.Error(1)
}

func (client *MockClient) StartContainer(id string, hostConfig *docker.HostConfig) error {

	args := client.Mock.Called(id, hostConfig)
	return args.Error(0)
}

func (client *MockClient) StopContainer(id string, timeout uint) error {

	args := client.Mock.Called(id, timeout)
	return args.Error(0)
}

func (client *MockClient) RestartContainer(id string, timeout uint) error {

	args := client.Mock.Called(id, timeout)
	return args.Error(0)
}

func (client *MockClient) PauseContainer(id string) error {

	args := client.Mock.Called(id)
	return args.Error(0)
}

func (client *MockClient) UnpauseContainer(id string) error {

	args := client.Mock.Called(id)
	return args.Error(0)
}

func (client *MockClient) TopContainer(id string, psArgs string) (docker.TopResult, error) {

	args := client.Mock.Called(id)
	return args.Get(0).(docker.TopResult), args.Error(1)
}

func (client *MockClient) Stats(opts docker.StatsOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) KillContainer(opts docker.KillContainerOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) RemoveContainer(opts docker.RemoveContainerOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) UploadToContainer(id string, opts docker.UploadToContainerOptions) error {

	args := client.Mock.Called(id, opts)
	return args.Error(0)
}

func (client *MockClient) DownloadFromContainer(id string, opts docker.DownloadFromContainerOptions) error {

	args := client.Mock.Called(id, opts)
	return args.Error(0)
}

func (client *MockClient) CopyFromContainer(opts docker.CopyFromContainerOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) WaitContainer(id string) (int, error) {

	args := client.Mock.Called(id)
	return args.Int(0), args.Error(1)
}

func (client *MockClient) CommitContainer(opts docker.CommitContainerOptions) (*docker.Image, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).(*docker.Image), args.Error(1)
}

func (client *MockClient) AttachToContainer(opts docker.AttachToContainerOptions) error {

	args := client.Mock.Called(opts)
	return args.Error(0)
}

func (client *MockClient) AttachToContainerNonBlocking(opts docker.AttachToContainerOptions) (docker.CloseWaiter, error) {

	args := client.Mock.Called(opts)

	return args.Get(0).(docker.CloseWaiter), args.Error(1)
}

func (client *MockClient) Logs(opts docker.LogsOptions) error {

	return errors.New("DockerClientWrapper imlementation of Logs currently doesn't allow to return a value")
}

func (client *MockClient) ResizeContainerTTY(id string, height, width int) error {

	args := client.Mock.Called(id, height, width)

	return args.Error(0)
}

func (client *MockClient) ExportContainer(opts docker.ExportContainerOptions) error {

	args := client.Mock.Called(opts)

	return args.Error(0)
}

func (client *MockClient) AddEventListener(listener chan <- *docker.APIEvents) error {

	args := client.Mock.Called(listener)

	return args.Error(0)
}

func (client *MockClient) RemoveEventListener(listener chan *docker.APIEvents) error {

	args := client.Mock.Called(listener)

	return args.Error(0)
}

func (client *MockClient) AuthCheck(conf *docker.AuthConfiguration) (docker.AuthStatus, error) {

	args := client.Mock.Called(conf)

	return args.Get(0).(docker.AuthStatus), args.Error(0)
}
