package config

import "github.com/joeshaw/envdecode"

var StorageDirectory = "./storage/app/"

type Config struct {
	AppName                  string   `env:"APP_NAME"`
	AppVersion               string   `env:"APP_VERSION"`
	AppEnv                   string   `env:"APP_ENV,default=development"`
	ApiHost                  string   `env:"API_HOST"`
	ApiRpcPort               string   `env:"API_RPC_PORT"`
	ApiPort                  string   `env:"API_PORT,default=8760"`
	ApiDocPort               uint16   `env:"API_DOC_PORT,default=8761"`
	ShutdownTimeout          uint     `env:"API_SHUTDOWN_TIMEOUT_SECONDS,default=30"`
	AllowedCredentialOrigins []string `env:"ALLOWED_CREDENTIAL_ORIGINS"`
	MiddlewareAddress        string   `env:"MIDDLEWARE_ADDR"`
	JwtExpireDaysCount       int      `env:"JWT_EXPIRE_DAYS_COUNT"`
	MysqlOption
	RabbitMQOption
	MongodbOption
	RedisOption
	PostgreSqlOption
}

// MysqlOption contains mySQL connection options
type MysqlOption struct {
	URI           string `env:"MYSQL_URI,default="`
	Pool          int    `env:"MYSQL_POOL,required"`
	SlowThreshold int    `env:"MYSQL_SLOW_LOG_THRESHOLD,required"`
}

type PostgreSqlOption struct {
	URI           string `env:"POSTGRE_URI,default="`
	Pool          int    `env:"POSTGRE_POOL,default=1000"`
	SlowThreshold int    `env:"POSTGRE_SLOW_LOG_THRESHOLD,default=200"`
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

type RedisOption struct {
	Host           string `env:"REDIS_HOST,required"`
	Password       string `env:"REDIS_PASSWORD"`
	ReadTimeoutMs  int16  `env:"REDIS_READ_TIMEOUT,required"`
	WriteTimeoutMs int16  `env:"REDIS_WRITE_TIMEOUT,required"`
}

func NewConfig() *Config {
	var cfg Config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
