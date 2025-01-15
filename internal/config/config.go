package config

type Config struct {
	RabbitMQ struct {
		URL      string
		Queue    string
		Exchange string
	}
	LogLevel string
}

// func LoadConfig() (*Config, error) {

// }
