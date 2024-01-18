package config

import (
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	r "github.com/ToraNoDora/little-sso/sso/internal/src/store/cache/redis_cache"
	p "github.com/ToraNoDora/little-sso/sso/internal/src/store/postgres"
)

type Config struct {
	Env      string        `yaml: "env" env-default:"locale"`
	TokenTtl time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GrpcConfig    `yaml:"grpc"`
	Store    p.Config
	Redis    r.Config
}

type GrpcConfig struct {
	Port    int           `yaml: "port"`
	Timeout time.Duration `yaml: "timeout"`
}

type envConfigs struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBMode     string `mapstructure:"DB_MODE"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisDB   int    `mapstructure:"REDIS_DB"`
	RedisPass string `mapstructure:"REDIS_PASSWORD"`
}

type ConfigsPaths struct {
	CfgPath string
	CfgName string
}

func MustLoad(configPath string) *Config {
	path := fetchConfigPath(configPath)
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(cfgPath string) *Config {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("file does not exist: " + cfgPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	secret := loadEnvVariables()
	cfg.Store = p.Config{
		Host:     secret.DBHost,
		Port:     secret.DBPort,
		Username: secret.DBUsername,
		Password: secret.DBPassword,
		DBName:   secret.DBName,
		SSLMode:  secret.DBMode,
	}
	cfg.Redis = r.Config{
		Host:     secret.RedisHost,
		Port:     secret.RedisPort,
		Password: secret.RedisPass,
		DB:       secret.RedisDB,
	}

	return &cfg
}

// priority: flag > env > default
func fetchConfigPath(cfgPath string) string {
	if cfgPath == "" {
		cfgPath = os.Getenv("CONFIG_PATH")
	}

	return cfgPath
}

// Call to load the variables from env
func loadEnvVariables() *envConfigs {
	loadEnv()

	return &envConfigs{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBMode:     os.Getenv("DB_MODE"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisDB:   getRedisDB(os.Getenv("REDIS_DB")),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
	}
}

const projectDirName = "sso"

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		panic("Error loading .env file")
	}
}

func getRedisDB(strVar string) int {
	n, err := strconv.Atoi(strVar)
	if err != nil {
		panic("error read redis db variables: " + err.Error())
	}

	return n
}
