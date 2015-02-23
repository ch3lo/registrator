package bridge

import (
	"log"
	"path"
	"strconv"
	"strings"
)

type Service struct {
	ID    string
	Name  string
	Port  int
	IP    string
	Tags  []string
	Attrs map[string]string
	TTL   int

	Origin ServicePort
}

func NewService(port ServicePort) *Service {
	s := new(Service)
	s.Origin = port
	s.setupMetadata()

	return s
}

func (s *Service) setId(hostname string) {
	s.ID = hostname + ":" + s.Origin.container.Name[1:] + ":" + s.Origin.ExposedPort

	if s.Origin.PortType == "udp" {
		s.ID = s.ID + ":udp"
	}

	id := mapDefault(s.Attrs, "id", "")
	if id != "" {
		s.ID = id
	}
}

func (s *Service) setName(isgroup bool) {
	container := s.Origin.container

	defaultName := strings.Split(path.Base(container.Config.Image), ":")[0]

	if isgroup {
		defaultName = defaultName + "-" + s.Origin.ExposedPort
	}

	s.Name = mapDefault(s.Attrs, "name", defaultName)
}

func (s *Service) setIp(internal bool) {
	if internal == true {
		s.IP = s.Origin.ExposedIP
	} else {
		s.IP = s.Origin.HostIP
	}
}

func (s *Service) setPort(internal bool) {
	var p int

	if internal == true {
		p, _ = strconv.Atoi(s.Origin.ExposedPort)
	} else {
		p, _ = strconv.Atoi(s.Origin.HostPort)
	}

	s.Port = p
}

func (s *Service) setTags(forceTags string) {
	if s.Origin.PortType == "udp" {
		s.Tags = combineTags(mapDefault(s.Attrs, "tags", ""), forceTags, "udp")
	} else {
		s.Tags = combineTags(mapDefault(s.Attrs, "tags", ""), forceTags)
	}
}

func (s *Service) ignore() bool {
	ignore := mapDefault(s.Attrs, "ignore", "")

	return ignore != ""
}

func (s *Service) setupMetadata() {
	container := s.Origin.container
	s.Attrs = serviceMetaData(container.Config.Env, s.Origin.ExposedPort)

	log.Println("service metadata:", s.Attrs)
}

func (s *Service) cleanMetadata() {
	delete(s.Attrs, "id")
	delete(s.Attrs, "tags")
	delete(s.Attrs, "name")
}
