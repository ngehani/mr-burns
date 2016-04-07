package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
)

type Client interface {
	Endpoint() string
	Ping() error
	ListVolumes(opts docker.ListVolumesOptions) ([]docker.Volume, error)
	CreateVolume(opts docker.CreateVolumeOptions) (*docker.Volume, error)
	InspectVolume(name string) (*docker.Volume, error)
	RemoveVolume(name string) error
	ListNetworks() ([]docker.Network, error)
	FilteredListNetworks(opts docker.NetworkFilterOpts) ([]docker.Network, error)
	NetworkInfo(id string) (*docker.Network, error)
	CreateNetwork(opts docker.CreateNetworkOptions) (*docker.Network, error)
	RemoveNetwork(id string) error
	ConnectNetwork(id string, opts docker.NetworkConnectionOptions) error
	DisconnectNetwork(id string, opts docker.NetworkConnectionOptions) error
	CreateExec(opts docker.CreateExecOptions) (*docker.Exec, error)
	StartExec(id string, opts docker.StartExecOptions) error
	StartExecNonBlocking(id string, opts docker.StartExecOptions) (docker.CloseWaiter, error)
	ResizeExecTTY(id string, height, width int) error
	InspectExec(id string) (*docker.ExecInspect, error)
	Version() (*docker.Env, error)
	Info() (*docker.DockerInfo, error)
	ListImages(opts docker.ListImagesOptions) ([]docker.APIImages, error)
	ImageHistory(name string) ([]docker.ImageHistory, error)
	RemoveImage(name string) error
	RemoveImageExtended(name string, opts docker.RemoveImageOptions) error
	InspectImage(name string) (*docker.Image, error)
	PushImage(opts docker.PushImageOptions, auth docker.AuthConfiguration) error
	PullImage(opts docker.PullImageOptions, auth docker.AuthConfiguration) error
	LoadImage(opts docker.LoadImageOptions) error
	ExportImage(opts docker.ExportImageOptions) error
	ExportImages(opts docker.ExportImagesOptions) error
	ImportImage(opts docker.ImportImageOptions) error
	BuildImage(opts docker.BuildImageOptions) error
	TagImage(name string, opts docker.TagImageOptions) error
	SearchImages(term string) ([]docker.APIImageSearch, error)
	SearchImagesEx(term string, auth docker.AuthConfiguration) ([]docker.APIImageSearch, error)
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
	UpdateContainer(id string, opts docker.UpdateContainerOptions) error
	RenameContainer(opts docker.RenameContainerOptions) error
	InspectContainer(id string) (*docker.Container, error)
	ContainerChanges(id string) ([]docker.Change, error)
	CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error)
	StartContainer(id string, hostConfig *docker.HostConfig) error
	StopContainer(id string, timeout uint) error
	RestartContainer(id string, timeout uint) error
	PauseContainer(id string) error
	UnpauseContainer(id string) error
	TopContainer(id string, psArgs string) (docker.TopResult, error)
	Stats(opts docker.StatsOptions) (retErr error)
	KillContainer(opts docker.KillContainerOptions) error
	RemoveContainer(opts docker.RemoveContainerOptions) error
	UploadToContainer(id string, opts docker.UploadToContainerOptions) error
	DownloadFromContainer(id string, opts docker.DownloadFromContainerOptions) error
	CopyFromContainer(opts docker.CopyFromContainerOptions) error
	WaitContainer(id string) (int, error)
	CommitContainer(opts docker.CommitContainerOptions) (*docker.Image, error)
	AttachToContainer(opts docker.AttachToContainerOptions) error
	AttachToContainerNonBlocking(opts docker.AttachToContainerOptions) (docker.CloseWaiter, error)
	Logs(opts docker.LogsOptions) error
	ResizeContainerTTY(id string, height, width int) error
	ExportContainer(opts docker.ExportContainerOptions) error
	AddEventListener(listener chan<- *docker.APIEvents) error
	RemoveEventListener(listener chan *docker.APIEvents) error
	AuthCheck(conf *docker.AuthConfiguration) error
}
