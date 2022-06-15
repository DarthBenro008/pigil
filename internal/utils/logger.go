package utils

import (
	"log"
	"os"
)

const (
	Red    = "\033[31m"
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

func ErrorLogger(err error, place string) {
	log.Printf("%spigil error (%s) : %s%s\n", Red, place, err.Error(), Reset)
	os.Exit(1)
}

func ErrorInformation(information string) {
	log.Printf("%spigil: %s%s\n", Red, information, Reset)
}

func InformationLogger(information string) {
	log.Printf("%spigil information: %s%s\n", Yellow, information, Reset)
}
