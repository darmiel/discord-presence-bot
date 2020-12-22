package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Token string
	Guild string
)

var (
	onlineCount uint16 = 0
	crawlStart  int64
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Guild, "g", "", "Guild-ID")
	flag.Parse()

	// .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(Token) == 0 {
		log.Println("Token not passed by parameter. Try to load from .env")
		Token = os.Getenv("DISCORD_TOKEN")
	}

	if len(Token) == 0 {
		log.Fatal("Discord token not found!")
	}
}

func main() {

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Error creating discord bot:", err)
		return
	}

	// Handlers
	discord.AddHandler(messageCreate)
	discord.AddHandler(onGuildMembersChunk)

	// Indents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error connecting:", err)
		return
	}

	// Check member ticker
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		checkMembers(discord)

		for {
			select {
			case t := <-ticker.C:
				log.Println("Ticker:", t)
				checkMembers(discord)
				log.Println()
				break
			}
		}
	}()
	//

	log.Println("Bot is not running! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINFO, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close
	err = discord.Close()
	if err != nil {
		log.Fatal("Error closing discord:", err)
	} else {
		log.Println("Bot closed!")
	}
}

/////////////////////////////////////////////////////////////////////////////////////

func checkMembers(s *discordgo.Session) {
	// update old count
	if onlineCount > 0 {
		log.Println("Online in the last", time.Now().Unix()-crawlStart, "seconds:", onlineCount)

		// save
		save()

		// reset current online members
		onlineCount = 0
	}
	crawlStart = time.Now().Unix()
	//

	log.Println("Checking if guild exists ...")
	guild, err := s.Guild(Guild)
	if err != nil {
		log.Println("Guild not found:", err)
		return
	}

	log.Println("Requesting members of guild:", guild.Name)
	err = s.RequestGuildMembers(Guild, "", 200, true)
	if err != nil {
		log.Fatal("Error requesting guild members:", err)
		return
	}
}

func isOnline(presence *discordgo.Presence) (res bool) {
	return presence.Status == "online" || presence.Status == "away" || presence.Status == "busy"
}

func getStats() (timespan uint32, online uint16) {
	unixNow := time.Now().Unix()
	return uint32(unixNow - crawlStart), onlineCount
}

func getStatsStr() (res string) {
	timespan, online := getStats()
	return fmt.Sprintf("Online between the last %v seconds: %v", timespan, online)
}

func save() {
	log.Println("(( saving ))")
}

/////////////////////////////////////////////////////////////////////////////////////

func onGuildMembersChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
	for _, presence := range c.Presences {
		if isOnline(presence) {
			onlineCount++
		}
	}

	now := time.Now().Unix()
	log.Println("Seconds:", now, "Online:", onlineCount)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore self
	if m.Author.ID == s.State.User.ID {
		log.Println("[Debug] Ignored message from bot")
		return
	}

	if m.Content == "!crawl" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Aight! Crawling ...", m.MessageReference)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		checkMembers(s)
	} else if m.Content == "!online" {
		message := getStatsStr()
		if _, err := s.ChannelMessageSendReply(m.ChannelID, message, m.MessageReference); err != nil {
			log.Println("Error:", err)
		}
	}
}
