package internal

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

/// Discord
var (
	Token       string
	Guild       string
	UpdateDelay time.Duration
)

/// Databae
var (
	SqlTable = "active_user_counts"
)

func InitFlags() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Guild, "g", "", "Guild-ID")
	flag.DurationVar(&UpdateDelay, "d", 30*time.Second, "Update delay")
	flag.Parse()
}

func InitDotEnv() {
	// .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	if len(Token) == 0 {
		log.Println("Token not passed by parameter. Trying to load from .env")

		if Token = os.Getenv("DISCORD_TOKEN"); len(Token) == 0 {
			log.Fatalln("Discord token not found!")
		}
	}
}
