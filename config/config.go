package config

import "os"

const (
	defaultGRPCHost       string = "localhost"
	defaultGRPCPort       string = "50051"
	defaultPrometheusHost string = "localhost"
	defaultPrometheusPort string = "9092"
	defaultDBHost         string = "localhost"
	defaultDBPort         string = "5432"
	defaultDBName         string = "postgres"
	defaultDBUser         string = "postgres"
	defaultMigrationsPath string = "migrations"
)

type Config struct {
	GRPCHost       string
	GRPCPort       string
	PrometheusHost string
	PrometheusPort string
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
		GRPCHost:       getEnv("GRPC_HOST", defaultGRPCHost),
		GRPCPort:       getEnv("GRPC_PORT", defaultGRPCPort),
		PrometheusHost: getEnv("PROMETHEUS_HOST", defaultPrometheusHost),
		PrometheusPort: getEnv("PROMETHEUS_PORT", defaultPrometheusPort),
		DBHost:         getEnv("DB_HOST", defaultDBHost),
		DBPort:         getEnv("DB_PORT", defaultDBPort),
		DBName:         getEnv("DB_NAME", defaultDBName),
		DBUser:         getEnv("DB_USER", defaultDBUser),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		MigrationsPath: getEnv("MIGRATIONS_PATH", defaultMigrationsPath),
	}
}
