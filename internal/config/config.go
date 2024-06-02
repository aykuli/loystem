package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

const (
	hostDefault = "localhost"
	portDefault = "8080"
)

var Options = Config{
	Address:              hostDefault + ":" + portDefault,
	DatabaseURI:          "",
	AccrualSystemAddress: "",
}

func init() {
	parseFlags()
}

func parseFlags() {
	fs := flag.NewFlagSet("loystem", flag.ContinueOnError)
	fs.StringVar(&Options.Address, "a", hostDefault+":"+portDefault, "server address to run on")
	fs.StringVar(&Options.DatabaseURI, "d", "", "database source name")
	fs.StringVar(&Options.AccrualSystemAddress, "r", "", "accrual system address")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Print(err)
	}
}
