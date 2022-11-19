package config

import "os"

const (
	defaultGRPCHost string = ""
	defaultGRPCPort string = "50051"
	defaultDBHost   string = "localhost"
	defaultDBPort   string = "5432"
	defaultDBName   string = "postgres"
	defaultDBUser   string = "postgres"
)

type Config struct {
	GRPCHost   string
	GRPCPort   string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
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
		GRPCHost:   getEnv("GRPC_HOST", defaultGRPCHost),
		GRPCPort:   getEnv("GRPC_PORT", defaultGRPCPort),
		DBHost:     getEnv("DB_HOST", defaultDBHost),
		DBPort:     getEnv("DB_PORT", defaultDBPort),
		DBName:     getEnv("DB_NAME", defaultDBName),
		DBUser:     getEnv("DB_USER", defaultDBUser),
		DBPassword: getEnv("DB_PASSWORD", ""),
	}
}
