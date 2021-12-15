package bx24sync

import "fmt"

type Socket struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s Socket) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Endpoint struct {
	Socket
	Method string
}

func (e Endpoint) URL(ssl bool) string {
	protocol := "http"
	if ssl {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s", protocol, e.Socket.String(), e.Method)
}

type BasicAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type RegistrarConfig struct {
	Http Socket
	ProducerConfig
}

type GeneratorConfig struct {
	DB              DataBaseConnection `json:"db"`
	Web             Socket             `json:"web"`
	StorageQueryTxt string             `json:"queryDir"`
}

type DataBaseConnection struct {
	Socket    `json:"socket"`
	BasicAuth `json:"auth"`
	Database  string `json:"database"`
}

type ProducerConfig struct {
	Broker     Socket `json:"broker"`
	Topic      string `json:"topic"`
}

type ConsumerConfig struct {
	Brokers    string `json:"brokers"`
	Topic      string `json:"topic"`
	GroupId    string `json:"groupId"`
	TargetLink string `json:"target"`
}

