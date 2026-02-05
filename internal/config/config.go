package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	RateLimitCapacity   int
	RateLimitRefillRate float64
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		RateLimitCapacity:   getEnvAsInt("RATE_LIMIT_CAPACITY", 10),
		RateLimitRefillRate: getEnvAsFloat("RATE_LIMIT_REFILL_RATE", 1.0),
	}
}

func getEnv(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if val, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if val, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.ParseFloat(val, 64); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}
