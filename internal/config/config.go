package config

import (
	"flag"
	"fmt"
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
	fmt.Println("-------\nPARSE FLAGS\n-------------\n")
	parseFlags()
}

func parseFlags() {
	fs := flag.NewFlagSet("loystem", flag.ContinueOnError)
	fs.StringVar(&Options.Address, "a", hostDefault+":"+portDefault, "server address to run on")
	fs.StringVar(&Options.DatabaseURI, "d", "", "database source name")
	fs.StringVar(&Options.AccrualSystemAddress, "r", "", "accrual system address")
	fmt.Println("-------\nOptions\n-------------\n", Options)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Print(err)
	}
}
