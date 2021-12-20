package bx24sync

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

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
	SSL    bool
}

func (e Endpoint) URL() string {
	protocol := "http"
	if e.SSL {
		protocol += "s"
	}

	if e.Method == "" {
		return fmt.Sprintf("%s://%s", protocol, e.String())
	} else {
		return fmt.Sprintf("%s://%s/%s", protocol, e.String(), e.Method)
	}
}

type BasicAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s BasicAuth) GetPair(sep string) string {
	return fmt.Sprintf("%s%s%s", s.User, sep, s.Password)
}

type RegistrarConfig struct {
	Http Socket
	ProducerConfig
	CheckInput bool
}

func NewRegistrarConfigFromEnv() RegistrarConfig {
	return RegistrarConfig{
		Http: Socket{
			Port: GetEnvAsInt("HTTP_PORT", 25410),
		},
		ProducerConfig: ProducerConfig{
			Broker: Socket{
				Host: getEnv("KAFKA_HOST", "127.0.0.1"),
				Port: GetEnvAsInt("KAFKA_PORT", 9092),
			},
			Topic: GetEnvAsString("KAFKA_TOPIC", ""),
		},
		CheckInput: GetEnvAsBool("HTTP_CHECK_INPUT", false),
	}
}

type GeneratorConfig struct {
	DB              DataBaseConnection `json:"db"`
	Web             Socket             `json:"web"`
	StorageQueryTxt string             `json:"queryDir"`
	CheckInput      bool               `json:"checkInput"`
}

func NewGeneratorConfigFromEnv() GeneratorConfig {
	return GeneratorConfig{
		DB: NewDataBaseConnectionFromEnv(),
		Web: Socket{
			Port: GetEnvAsInt("HTTP_PORT", 8080),
		},
		StorageQueryTxt: "./sql",
		CheckInput:      GetEnvAsBool("HTTP_CHECK_INPUT", true),
	}
}

type DataBaseConnection struct {
	Socket    `json:"socket"`
	BasicAuth `json:"auth"`
}

func NewDataBaseConnectionFromEnv() DataBaseConnection {
	return DataBaseConnection{
		Socket: Socket{
			Host: GetEnvAsString("DB_HOST", "127.0.0.1"),
			Port: GetEnvAsInt("DB_PORT", 1433),
		},
		BasicAuth: BasicAuth{
			User:     GetEnvAsString("DB_USER", ""),
			Password: GetEnvAsString("DB_PASSWORD", ""),
		},
	}
}

func (c DataBaseConnection) MakeConnURL() *url.URL {
	return &url.URL{
		Scheme: "sqlserver",
		Host:   c.Socket.String(),
		User:   url.UserPassword(c.User, c.Password),
	}
}

type ProducerConfig struct {
	Broker Socket `json:"broker"`
	Topic  string `json:"topic"`
}

type ConsumerConfig struct {
	Brokers           []Socket `json:"brokers"`
	Topic             string   `json:"topic"`
	GroupId           string   `json:"groupId"`
	Partition         int      `json:"Partition"`
	GeneratorEndpoint Endpoint `json:"generatorEndpoint"`
	TargetEndpoint    Endpoint `json:"targetEndpoint"`
}

func NewConsumerConfigFromEnv() (ConsumerConfig, error) {
	c := ConsumerConfig{}

	brokers := GetEnvAsStringSlice("KAFKA_BROKERS", ",", make([]string, 0))

	for _, v := range brokers {
		v := strings.ReplaceAll(v, " ", "")
		socketPath := strings.Split(v, ":")

		if len(socketPath) < 2 {
			return c, fmt.Errorf("don't manage to parse socket from string. Raw %s", v)
		}

		host := socketPath[0]
		port, err := strconv.Atoi(socketPath[1])

		if err != nil {
			return c, fmt.Errorf("don't manage to parse int from string. Raw %s. Error %+v", socketPath[1], err)
		}

		c.Brokers = append(c.Brokers, Socket{
			Host: host,
			Port: port,
		})

		c.Topic = GetEnvAsString("KAFKA_TOPIC", "")
		c.GroupId = GetEnvAsString("KAFKA_GROUP_ID", "")
		c.Partition = GetEnvAsInt("KAFKA_PARTITION", 0)

		c.GeneratorEndpoint = Endpoint{
			Socket: Socket{
				Host: GetEnvAsString("GENERATOR_HOST", ""),
				Port: GetEnvAsInt("GENERATOR_PORT", 0),
			},
			SSL: GetEnvAsBool("GENERATOR_SSL", false),
		}

		c.TargetEndpoint = Endpoint{
			Socket: Socket{
				Host: GetEnvAsString("TARGET_HOST", ""),
				Port: GetEnvAsInt("TARGET_PORT", 0),
			},
			SSL: GetEnvAsBool("TARGET_SSL", false),
		}
	}

	return c, nil
}