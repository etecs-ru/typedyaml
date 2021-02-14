package example

import "time"

//go:generate go run github.com/etecs-ru/typedyaml -interface Gateway UserGateway OrdersGateway

// Gateway acts like a polymorphic type that could be
// marshalled to yaml as well as unmarshalled from the format.
type Gateway interface {
	TypedYAML(*GatewayTyped) string
}

// MicroserviceConfig is a sample microservice configuration that needs to describe several
// instances of some entity that has different kinds.
type MicroserviceConfig struct {
	Tags     map[string]string `yaml:"tags"`
	Gateway  GatewayTyped      `yaml:"gateway"`
	Gateways []GatewayTyped    `yaml:"gateways"`
}

// UserGateway acts as an example of Gateway of kind User
type UserGateway struct {
	Enabled            bool   `yaml:"enabled"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	User               string `yaml:"user"`
	SecurityDescriptor string `yaml:"security_descriptor"`
}

// OrdersGateway acts as an example of Gateway of kind Orders
type OrdersGateway struct {
	Enabled bool               `yaml:"enabled"`
	Kafka   KafkaConfiguration `yaml:"kafka"`
}

// KafkaConfiguration is an additional "weight" to an Orders struct, just to ensure in tests
// that this tool could work in a real-world scenarios.
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
