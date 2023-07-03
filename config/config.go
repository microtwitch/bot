package config

import (
	"log"
	"os"
)

var Channel string
var Admin string
var Token string
var BotUser string
var ReceiverTarget string
var EdgeTarget string

func Init() {
	channel, present := os.LookupEnv("CHANNEL")
	if !present {
		log.Fatalln("No channel found")
	} else {
		Channel = channel
	}

	admin, present := os.LookupEnv("ADMIN")
	if !present {
		log.Fatalln("No admin found")
	} else {
		Admin = admin
	}

	token, present := os.LookupEnv("TOKEN")
	if !present {
		log.Fatalln("No token found")
	} else {
		Token = token
	}

	botUser, present := os.LookupEnv("BOT_USER")
	if !present {
		log.Fatalln("No bot user found ")
	} else {
		BotUser = botUser
	}

	receiverTarget, present := os.LookupEnv("RECEIVER_TARGET")
	if !present {
		log.Fatalln("No receiver target found ")
	} else {
		ReceiverTarget = receiverTarget
	}

	edgeTarget, present := os.LookupEnv("EDGE_TARGET")
	if !present {
		log.Fatalln("No edge target found ")
	} else {
		EdgeTarget = edgeTarget
	}
}
