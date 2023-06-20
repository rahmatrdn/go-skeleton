package config

import "github.com/joeshaw/envdecode"

var StorageDirectory = "./storage/app/"

type Config struct {
	AppEnv                   string   `env:"APP_ENV,default=development"`
	ApiHost                  string   `env:"API_HOST"`
	ApiPort                  string   `env:"API_PORT,default=8760"`
	ApiDocPort               uint16   `env:"API_DOC_PORT,default=8761"`
	ShutdownTimeout          uint     `env:"API_SHUTDOWN_TIMEOUT_SECONDS,default=30"`
	AllowedCredentialOrigins []string `env:"ALLOWED_CREDENTIAL_ORIGINS"`
	MiddlewareAddress        string   `env:"MIDDLEWARE_ADDR"`
	JwtSecret                string   `env:"JWT_SECRET_KEY"`
	JwtExpireDaysCount       int      `env:"JWT_EXPIRE_DAYS_COUNT"`
	MysqlOption
	RedisOption
	RabbitMQOption
	MongodbOption
	TelegramBotOption
}

// MysqlOption contains mySQL connection options
type MysqlOption struct {
	Driver       string `env:"MYSQL_DRIVER,default="`
	Host         string `env:"MYSQL_HOST,required"`
	Port         string `env:"MYSQL_PORT,required"`
	Pool         int    `env:"MYSQL_POOL,required"`
	DatabaseName string `env:"MYSQL_DATABASE_NAME,required"`
	Username     string `env:"MYSQL_USERNAME,required"`
	Password     string `env:"MYSQL_PASSWORD"`
	TimeZone     string `env:"MYSQL_TIMEZONE,required"`
}

type RedisOption struct {
	Host           string `env:"REDIS_HOST,required"`
	Password       string `env:"REDIS_PASSWORD"`
	ReadTimeoutMs  int    `env:"REDIS_READ_TIMEOUT,default=300"`
	WriteTimeoutMs int    `env:"REDIS_WRITE_TIMEOUT,default=300"`
}

type RabbitMQOption struct {
	Uri             string `env:"RABBITMQ_URI,required"`
	Exchange        string `env:"RABBITMQ_EXCHANGE,default=events"`
	QueueType       string `env:"RABBITMQ_QUEUE_TYPE,default=topic"`
	QueuePrefix     string `env:"RABBITMQ_QUEUE_PREFIX,default=Ngorder API"`
	QueueRetryCount int    `env:"RABBITMQ_RETRY_COUNT,default=3"`
}

type MongodbOption struct {
	Uri          string `env:"MONGODB_URI,required"`
	DatabaseName string `env:"MONGODB_DATABASE_NAME,required"`
}

type TelegramBotOption struct {
	Url    string `env:"TELE_URL"`
	Token  string `env:"TOKENTELE"`
	ChatID string `env:"CHATIDTELE"`
}

func NewConfig() *Config {
	var cfg Config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
