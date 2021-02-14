package example

import "time"

//go:generate go run github.com/etecs-ru/typedyaml -interface Gateway UserGateway OrdersGateway
type Gateway interface {
	TypedYAML(*GatewayTyped) string
}

type MicroserviceConfig struct {
	Tags     map[string]string `yaml:"tags"`
	Gateway  GatewayTyped      `yaml:"gateway"`
	Gateways []GatewayTyped    `yaml:"gateways"`
}

type UserGateway struct {
	Enabled            bool   `yaml:"enabled"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	User               string `yaml:"user"`
	SecurityDescriptor string `yaml:"security_descriptor"`
}

type OrdersGateway struct {
	Enabled bool               `yaml:"enabled"`
	Kafka   KafkaConfiguration `yaml:"kafka"`
}

type KafkaConfiguration struct {
	Hosts          []string      `yaml:"hosts" validate:"required"`
	Topics         []string      `yaml:"topics" validate:"required"`
	GroupID        string        `yaml:"group_id" validate:"required"`
	ClientID       string        `yaml:"client_id"`
	ConnectBackoff time.Duration `yaml:"connect_backoff" validate:"min=0"`
	ConsumeBackoff time.Duration `yaml:"consume_backoff" validate:"min=0"`
	WaitClose      time.Duration `yaml:"wait_close" validate:"min=0"`
	MaxWaitTime    time.Duration `yaml:"max_wait_time"`
	IsolationLevel int           `yaml:"isolation_level"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
}
