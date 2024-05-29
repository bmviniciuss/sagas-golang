package env

import "github.com/caarlos0/env"

type (
	Config struct {
		App   App
		Kafka Kafka
	}

	App struct {
		ServiceName string
	}

	Kafka struct {
		BootstrapServers string
		Topics           []string
		GroupID          string
	}
)

type config struct {
	ServiceName           string `env:"SERVICE_NAME" envDefault:"customers"`
	KafkaBootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"localhost:9092"`
	KafkaTopics           string `env:"KAFKA_TOPICS" envDefault:"service.customers.request"`
	KafkaGroupID          string `env:"KAFKA_GROUP_ID" envDefault:"customer-service-group"`
}

func Load() (*config, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
