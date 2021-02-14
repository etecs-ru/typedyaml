# typedyaml

[![Go Reference](https://pkg.go.dev/badge/github.com/etecs-ru/typedyaml.svg)](https://pkg.go.dev/github.com/etecs-ru/typedyaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/etecs-ru/typedyaml)](https://goreportcard.com/report/github.com/etecs-ru/typedyaml)

This is a code generator for Go based on [github.com/etecs-ru/typedjson](https://github.com/etecs-ru/typedjson) that alleviates YAML marshaling/unmarshalling unrelated structs in typed fashion.

Tool depends on [github.com/goccy/go-yaml](https://github.com/goccy/go-yaml) allowing such benefits as [fields validation](https://github.com/goccy/go-yaml/blob/868d322819b933bce2a46cfa2951c08706600f14/validate_test.go#L75).

Imagine, that you need to configure several instances of some service, each of different kind
using a YAML object with some key `Config`. 
The value of this field can correspond to two structs of your Go program: `FooConfig` and `BarConfig`.
So, the field `Config` in your struct must be able to hold a value of two possible types.
In this case, you have the following options:

1. You can declare field `Config` as `interface{}` and somehow determine what type you should expect, assign an object of this type to `Config` 
and then unmarshal object.
1. You can unmarshal field `Config` separately.
1. You can implement custom `MarshalYAML`/`UnmarshalYAML` for the third type that automatically will handle these cases.

This package provides means to generate all boilerplate code for the third case.

## Usage

```sh
typedyaml [OPTION] NAME...
```

Options:

* `-interface` string

	Name of the interface that encompass all types.

* `-output` string

	Output path where generated code should be saved.

* `-package` string

	Package name in generated file (default to GOPACKAGE).

* `-typed` string

	The name of the struct that will be used for typed interface (default to `{{interface}}{{Typed}}`).

Each name in position argument should be the name of the struct. 
You can set an alias for struct name like this: `foo=*FooConfig`.

## Example

See code in the [/example](https://github.com/etecs-ru/typedyaml/tree/master/example) folder.

For example, you have some microservice that serves orders to users and use gateways to other microservices that provides it. You need to configure it in one YAML file in a smart way.

You define `UserGateway` and `OrdersGateway` structures in our Go code, like:

```go
package config

import "time"

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
```

So you need to group this pieces together in some `MicroserviceConfig` structure and want to read it from YAML file.

We propose to implement polymorphic type `Gateway` using this code generation tool.

First, you must declare an interface that will hold either of these structs. The interface must have the method `TypedYAML` with a signature holding name of your container struct with a `Typed` suffix, like *Gateway:GatewayTyped*. This method will advise the compiler to work with types.

Let's see:

```go
package config 

//go:generate go run github.com/etecs-ru/typedyaml -package config -interface Gateway UserGateway OrdersGateway
type Gateway interface {
	TypedYAML(*GatewayTyped) string
}
```
Now, run `go generate`.
Generated struct `ConfigTyped` will have special implemented methods `MarshalYAML`/`UnmarshalYAML`. `GatewayTyped` could be used as a single instance, or in a slice. Adding configuration for new gateways now working like a charm -- just add its type to code generation argument and regenerate the code.

Let us write some configuration example:

```yaml
tags:
  key: value
gateway:
  type: UserGateway
  value:
	enabled: false
	host: internal-users.microservice.lan
	port: 8443
	user: robot
	security_descriptor: SY
gateways:
- type: UserGateway
  value:
	enabled: true
	host: external-users.microservice.lan
	port: 8443
	user: robot
	security_descriptor: SY
- type: OrdersGateway
  value:
	enabled: true
	kafka:
	  hosts:
	  - host-a:7891
	  - host-b:7892
	  topics:
	  - topic-a
	  - topic-b
	  - topic-c
	  group_id: gid01
	  client_id: cid01
	  connect_backoff: 1000000000
	  consume_backoff: 1000000000
	  wait_close: 1000000000
	  max_wait_time: 1000000000
	  isolation_level: 0
	  username: user
	  password: password
```


You can use generated code like this:

```go
package config

import (
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

type MicroserviceConfig struct {
	Tags     map[string]string `yaml:"tags"`
	Gateway  GatewayTyped      `yaml:"gateway"`
	Gateways []GatewayTyped    `yaml:"gateways"`
}

func ReadMicroserviceConfigFromFile(path string) (*MicroserviceConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg MicroserviceConfig
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func WriteMicroserviceConfigToFile(cfg MicroserviceConfig, path string) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

```
