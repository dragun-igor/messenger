package config

import "os"

const (
	defaultGRPCAddr       string = "localhost:50051"
	defaultPrometheusAddr string = "localhost:9092"
	defaultDBHost         string = "localhost"
	defaultDBPort         string = "5432"
	defaultDBName         string = "postgres"
	defaultDBUser         string = "postgres"
	defaultMigrationsPath string = "migrations"
)

type Config struct {
	GRPCAddr       string
	PrometheusAddr string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	MigrationsPath string
}

func getEnv(key, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return env
}

func Get() *Config {
	return &Config{
		GRPCAddr:       getEnv("GRPC_ADDR", defaultGRPCAddr),
		PrometheusAddr: getEnv("PROMETHEUS_ADDR", defaultPrometheusAddr),
		DBHost:         getEnv("DB_HOST", defaultDBHost),
		DBPort:         getEnv("DB_PORT", defaultDBPort),
		DBName:         getEnv("DB_NAME", defaultDBName),
		DBUser:         getEnv("DB_USER", defaultDBUser),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		MigrationsPath: getEnv("MIGRATIONS_PATH", defaultMigrationsPath),
	}
}
