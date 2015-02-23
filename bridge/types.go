//go:generate go-extpoints . AdapterFactory
package bridge

import (
	"net/url"

	dockerapi "github.com/fsouza/go-dockerclient"
)

type AdapterFactory interface {
	New(uri *url.URL) RegistryAdapter
}

type RegistryAdapter interface {
	Ping() error
	Register(service *Service) error
	Deregister(service *Service) error
	Refresh(service *Service) error
}

type Config struct {
	HostIp          string
	Internal        bool
	ForceTags       string
	RefreshTtl      int
	RefreshInterval int
}

type ServicePort struct {
	HostPort          string
	HostIP            string
	ExposedPort       string
	ExposedIP         string
	PortType          string
	ContainerHostname string
	ContainerID       string
	container         *dockerapi.Container
}
