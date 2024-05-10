package env

import "github.com/caarlos0/env"

type config struct {
	ServiceName           string `env:"SERVICE_NAME" envDefault:"accounting"`
	KafkaBootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"localhost:9092"`
	KafkaTopics           string `env:"KAFKA_TOPICS" envDefault:"service.accounting.request"`
	KafkaGroupID          string `env:"KAFKA_GROUP_ID" envDefault:"accounting-service-group"`
}

func Load() (*config, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
