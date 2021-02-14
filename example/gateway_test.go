package example

import (
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestGateway(t *testing.T) {
	original := MicroserviceConfig{
		Tags: map[string]string{"key": "value"},
		Gateway: GatewayTyped{
			Gateway: &UserGateway{
				Enabled:            false,
				Host:               "internal-users.microservice.lan",
				Port:               8443,
				User:               "robot",
				SecurityDescriptor: "SY",
			},
		},
		Gateways: []GatewayTyped{
			{
				Gateway: &UserGateway{
					Enabled:            true,
					Host:               "external-users.microservice.lan",
					Port:               8443,
					User:               "robot",
					SecurityDescriptor: "SY",
				},
			},
			{
				Gateway: &OrdersGateway{
					Enabled: true,
					Kafka: KafkaConfiguration{
						Hosts:          []string{"host-a:7891", "host-b:7892"},
						Topics:         []string{"topic-a", "topic-b", "topic-c"},
						GroupID:        "gid01",
						ClientID:       "cid01",
						ConnectBackoff: time.Second,
						ConsumeBackoff: time.Second,
						WaitClose:      time.Second,
						MaxWaitTime:    time.Second,
						IsolationLevel: 0,
						Username:       "user",
						Password:       "password",
					},
				},
			},
		},
	}
	b, err := yaml.Marshal(original)
	assert.NoError(t, err)

	var control MicroserviceConfig

	err = yaml.Unmarshal(b, &control)
	assert.NoError(t, err)
	assert.Equal(t, original, control)
}
