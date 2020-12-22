package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
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

	// Start "timer"
	checkMembers(discord)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore self
	if m.Author.ID == s.State.User.ID {
		fmt.Println("[Debug] Ignored message from bot")
		return
	}

	if m.Content == "!test" {
		message, err := s.ChannelMessageSendReply(m.ChannelID, "Test angekommen!", m.MessageReference)
		if err != nil {
			fmt.Println("Error sending message:", err)
		} else {
			fmt.Println("Message sent:", message)
		}
	}
}

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

func onGuildMembersChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
	fmt.Println("Got chunk:")
	fmt.Println(c)

	for i, member := range c.Members {
		presence := c.Presences[i]
		fmt.Println("Found user in chunk:", member.User.Username, "<->", presence.Status)
	}
}
