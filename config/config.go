package config

import (
	"log"
	"os"
)

var Channel string

func Init() {
	channel, present := os.LookupEnv("CHANNEL")
	if !present {
		log.Fatalln("No channel found")
	} else {
		Channel = channel
	}
}
