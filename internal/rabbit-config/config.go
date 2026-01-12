package rabbitconfig

type RabbitMQConfig struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
}

func GetRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URL:        "amqp://guest:guest@localhost:5672/",
		Exchange:   "logs",
		Queue:      "app_logs",
		RoutingKey: "log.info",
	}
}
