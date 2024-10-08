package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agent/handlers"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agent/hasher"
)

var (
	flagNetAddr    string
	reportInterval *int
	pollInterval   *int
	flagLogLevel   string
	flagKey        string
	rateLimit      *int
)

func parseFlags() {
	flag.StringVar(&flagNetAddr, "a", ":8080", "address and port to run server")

	reportInterval = flag.Int("r", 10, "report interval")
	pollInterval = flag.Int("p", 2, "poll interval")
	flag.StringVar(&flagLogLevel, "log", "info", "log level")
	flag.StringVar(&flagKey, "k", "", "key for hashing data")
	rateLimit = flag.Int("l", 1, "count of concurrent messages to server")

	flag.Parse()

	// для случаев, когда в переменной окружения ADDRESS присутствует непустое значение,
	// переопределим адрес агента,
	// даже если он был передан через аргумент командной строки
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagNetAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		val, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatalln("Environment variable \"REPORT_INTERVAL\" must be int")
		}
		*reportInterval = val
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		val, err := strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatalln("Environment variable \"POLL_INTERVAL\" must be int")
		}
		*pollInterval = val
	}
	if envLogLevel := os.Getenv("AGENT_LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		flagKey = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		val, err := strconv.Atoi(envRateLimit)
		if err != nil {
			log.Fatalln("Environment variable \"POLL_INTERVAL\" must be int")
		}
		*rateLimit = val
	}

	handlers.SetReportInterval(time.Duration(*reportInterval))
	handlers.SetPollInterval(time.Duration(*pollInterval))
	hasher.SetKey(flagKey)
}
