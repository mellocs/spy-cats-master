package config

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Env        string
	HTTPServer HTTPServer
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func MustLoad() *Config {
	err := loadEnv("./config/.env")
	if err != nil {
		log.Fatal(err)
	}

	timeoutStr := os.Getenv("TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Print("CONFIG_TIMEOUT could not be parsed")
	}

	httpServer := HTTPServer{
		Address: os.Getenv("HOST_ADDRESS"),
		Timeout: timeout,
	}

	return &Config{
		Env:        os.Getenv("ENV"),
		HTTPServer: httpServer,
	}
}

func loadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		setErr := os.Setenv(key, value)
		if setErr != nil {
			return err
		}
	}

	return scanner.Err()
}
