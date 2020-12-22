package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Token string
	Guild string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Guild, "g", "", "Guild-ID")
	flag.Parse()
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Handlers
	discord.AddHandler(messageCreate)
	discord.AddHandler(onGuildMembersChunk)

	// Indents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	// Check member ticker
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		select {
		case <-ticker.C:
			checkMembers(discord)
			break
		}
	}()
	checkMembers(discord)
	//

	fmt.Println("Bot is not running! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINFO, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close
	err = discord.Close()
	if err != nil {
		fmt.Println("Error closing discord:", err)
	} else {
		fmt.Println("Bot closed!")
	}
}

/////////////////////////////////////////////////////////////////////////////////////

func checkMembers(s *discordgo.Session) {
	guild, err := s.Guild(Guild)
	if err != nil {
		fmt.Println("Guild not found!")
		return
	}

	fmt.Println("Requesting members of guild:", guild.Name)

	err = s.RequestGuildMembers(Guild, "", 200, true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

/////////////////////////////////////////////////////////////////////////////////////

func onGuildMembersChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
	fmt.Println("Got chunk:")
	fmt.Println(c)

	for i, member := range c.Members {
		presence := c.Presences[i]
		fmt.Println("Found user in chunk:", member.User.Username, "<->", presence.Status)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore self
	if m.Author.ID == s.State.User.ID {
		fmt.Println("[Debug] Ignored message from bot")
		return
	}

	if m.Content == "!crawl" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Aight! Crawling ...", m.MessageReference)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		checkMembers(s)
	}
}
