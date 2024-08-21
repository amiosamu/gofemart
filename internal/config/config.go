package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port              string
	DBPort            string
	ScoringSystemPort string
	TokenTTL          time.Duration
	LogLevel          string
}

func NewConfig() *Config {
	return &Config{
		Port:     ":8080",
		TokenTTL: time.Minute * 30,
		LogLevel: "debug",
	}
}

type NetAddr struct {
	Host string
	Port int
}

func (a NetAddr) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddr) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

func (c *Config) ParseFlags() {
	port := new(NetAddr)
	_ = flag.Value(port)

	flag.Var(port, "a", "net address host:port")
	dbPort := flag.String("d", "", "port for database")
	scoringSystemPort := flag.String("r", "", "port for scoring system")

	flag.Parse()
	c.DBPort = *dbPort
	c.ScoringSystemPort = *scoringSystemPort

	if port.String() != ":0" {
		c.Port = port.String()
	}

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		c.Port = envRunAddr
	}

	if envPath := os.Getenv("DATABASE_URI"); envPath != "" {
		c.DBPort = envPath
	}

	if envScoring := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envScoring != "" {
		c.ScoringSystemPort = envScoring
	}

}
