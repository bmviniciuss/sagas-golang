package env

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
