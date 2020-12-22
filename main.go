package main

import (
	"database/sql"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/// Discord
var (
	Token       string
	Guild       string
	UpdateDelay time.Duration
)

/// Database
const table = "active_user_counts"

var db *sql.DB

/// Counting
var (
	onlineCounts = make(map[string]map[discordgo.Status]uint16)
	crawlStarts  = make(map[string]int64)
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Guild, "g", "", "Guild-ID")
	flag.DurationVar(&UpdateDelay, "d", 30000, "Update delay")
	flag.Parse()

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

	initDatabase()
}

func main() {

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("Error creating discord bot:", err)
		return
	}

	// Handlers
	discord.AddHandler(messageCreate)
	discord.AddHandler(onGuildMembersChunk)

	// Indents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if err := discord.Open(); err != nil {
		log.Fatalln("Error connecting:", err)
		return
	}

	// Check member ticker
	log.Println("Update delay:", UpdateDelay)
	ticker := time.NewTicker(UpdateDelay * time.Millisecond)
	go func() {
		checkMembers(discord)

		for {
			select {
			case <-ticker.C:
				checkMembers(discord)
				break
			}
		}
	}()
	//

	// Close
	defer func() {
		log.Println("Closing bot ...")
		if err := discord.Close(); err != nil {
			log.Fatalln("Error closing discord:", err)
		} else {
			log.Println("Bot closed!")
		}
	}()

	log.Println("Bot is not running! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINFO, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
