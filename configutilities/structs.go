package configutilities

type DatabasePostgresConfiguration struct {
	Host         string `env:"DATABASE_HOST"     yaml:"host"     validate:"required"`
	Port         int    `env:"DATABASE_PORT"     yaml:"port"     validate:"inhouseports"`
	User         string `env:"DATABASE_USERNAME" yaml:"username" validate:"required"`
	Pass         string `env:"DATABASE_PASS"     yaml:"pass"     validate:"required"`
	DatabaseName string `env:"DATABASE_NAME"     yaml:"name"     validate:"required"`
	SSLMode      string `env:"DATABASE_SSL_MODE" yaml:"ssl_mode" validate:"sslmode"`
}

type SlogConfiguration struct {
	Level string `env:"SLOG_LEVEL" yaml:"level" validate:"loglevel"`
}

type GinConfiguration struct {
	Mode string `env:"GIN_MODE" yaml:"mode" validate:"ginmode"`
	Port int    `env:"APP_PORT" yaml:"port" validate:"inhouseports"`
}

type AWSSDKConfiguration struct {
	EndpointURL     string `env:"AWS_ENDPOINT_URL"      yaml:"endpointURL"`
	Region          string `env:"AWS_REGION"            yaml:"region"`
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"     yaml:"accessKeyID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" yaml:"secretAccessKey" sensitive:"true"`
}
