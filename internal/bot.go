package internal

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

/// Counting
var (
	OnlineCounts = make(map[string]map[discordgo.Status]uint16)
	CrawlStarts  = make(map[string]int64)
)

func CreateBot() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("Error creating discord bot:", err)
		return
	}

	// Handlers
	discord.AddHandler(MessageCreate)
	discord.AddHandler(GuildMembersChunk)

	// Indents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if err := discord.Open(); err != nil {
		log.Fatalln("Error connecting:", err)
		return
	}

	// Check member ticker
	log.Println("Update delay:", UpdateDelay)
	ticker := time.NewTicker(UpdateDelay)
	go func() {
		CheckMembers(discord)

		for {
			select {
			case <-ticker.C:
				CheckMembers(discord)
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

	// TODO: Change me
	c := make(chan bool)
	<-c
}
