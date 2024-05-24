package config

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	RequestMaxRetries    int
}

const (
	hostDefault = "localhost"
	portDefault = "8080"
)

var Options = Config{
	Address:              hostDefault + ":" + portDefault,
	DatabaseURI:          "",
	AccrualSystemAddress: "",
	RequestMaxRetries:    3,
}

func init() {
	parseFlags()
	err := env.Parse(&Options)
	if err != nil {
		log.Fatal(err)
	}

	if Options.DatabaseURI == "" {
		log.Fatal("Error: database URI is required")
	}
	if Options.AccrualSystemAddress == "" {
		log.Fatal("Error: accrual system address is required")
	}
}

func parseFlags() {
	fs := flag.NewFlagSet("loystem", flag.ContinueOnError)
	fs.StringVar(&Options.Address, "a", hostDefault+":"+portDefault, "server address to run on")
	fs.StringVar(&Options.DatabaseURI, "d", "", "database source name")
	fs.StringVar(&Options.AccrualSystemAddress, "r", "", "accrual system address")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
