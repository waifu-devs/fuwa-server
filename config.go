package main

import (
	"log/slog"
	"os"
	"path"
	"strconv"
)

type config struct {
	port             string
	dataPath         string
	logLevel         slog.Level
	storageHost      string
	storageRegion    string
	storageAccessKey string
	storageSecretKey string
	isWorker         bool
	encryptionKey    []byte
	signingKey       []byte
	otelEndpoint     string
	otelServiceName  string
	otelSamplingRate float64
	otelAuthToken    string
	otelLogEndpoint  string
}

func loadConfigFromEnv() *config {
	c := &config{
		port:             os.Getenv("PORT"),
		dataPath:         os.Getenv("FUWA_DATA"),
		logLevel:         slog.LevelInfo,
		storageHost:      os.Getenv("FUWA_STORAGE_HOST"),
		storageRegion:    os.Getenv("FUWA_STORAGE_REGION"),
		storageAccessKey: os.Getenv("FUWA_STORAGE_ACCESS_KEY"),
		storageSecretKey: os.Getenv("FUWA_STORAGE_SECRET_KEY"),
		encryptionKey:    []byte(os.Getenv("FUWA_ENCRYPTION_KEY")),
		signingKey:       []byte(os.Getenv("FUWA_SIGNING_KEY")),
		otelEndpoint:     os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		otelServiceName:  os.Getenv("OTEL_SERVICE_NAME"),
		otelSamplingRate: 1.0,
		otelAuthToken:    os.Getenv("OTEL_AUTH_TOKEN"),
		otelLogEndpoint:  os.Getenv("OTEL_LOG_ENDPOINT"),
	}
	if c.port == "" {
		c.port = "6942"
	}
	if c.dataPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		c.dataPath = path.Join(homeDir, ".fuwa")
	} else {
		c.dataPath = path.Join(c.dataPath, ".fuwa")
	}
	switch os.Getenv("FUWA_LOG_LEVEL") {
	case "debug":
		c.logLevel = slog.LevelDebug
	}
	switch os.Getenv("FUWA_IS_WORKER") {
	case "1", "true":
		c.isWorker = true
	}
	if rate := os.Getenv("OTEL_SAMPLING_RATE"); rate != "" {
		if parsedRate, err := strconv.ParseFloat(rate, 64); err == nil {
			c.otelSamplingRate = parsedRate
		}
	}
	return c
}
