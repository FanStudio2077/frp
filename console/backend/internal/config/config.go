package config

import "os"

// Config is the console runtime configuration, loaded from environment variables.
type Config struct {
	ListenAddr string
	DBPath     string
	JWTSecret  string
	MasterKey  string
}

// Load reads configuration from the environment with sensible defaults.
func Load() Config {
	return Config{
		ListenAddr: getenv("CONSOLE_LISTEN", ":8080"),
		DBPath:     getenv("CONSOLE_DB", "./console.db"),
		JWTSecret:  getenv("CONSOLE_JWT_SECRET", "change-me-jwt-secret"),
		MasterKey:  getenv("CONSOLE_MASTER_KEY", "change-me-32-byte-master-key!!"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
