package bx24sync

import "fmt"

type Socket struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type GeneratorConfig struct {
	DB              DataBaseAuth `json:"db"`
	Web             Socket       `json:"web"`
	StorageQueryTxt string       `json:"queryDir"`
}

type Auth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type DataBaseAuth struct {
	Socket   `json:"socket"`
	Auth     `json:"auth"`
	Database string `json:"database"`
}

type RegistrarConfig struct {
	ProducerConfig
	Endpoint
}

type ProducerConfig struct {
	Broker     Socket `json:"broker"`
	Topic      string `json:"topic"`
	Partition  int    `json:"partition"`
	MessageKey string `json:"messageKey"`
}

type Endpoint struct {
	Socket
	Method string
}

type ConsumerConfig struct {
	Brokers    string `json:"brokers"`
	Topic      string `json:"topic"`
	GroupId    string `json:"groupId"`
	TargetLink string `json:"target"`
}

func (s Socket) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
