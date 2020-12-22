package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

func onGuildMembersChunk(_ *discordgo.Session, c *discordgo.GuildMembersChunk) {
	for _, presence := range c.Presences {
		onlineCount, ok := onlineCounts[c.GuildID]
		if !ok {
			onlineCount = make(map[discordgo.Status]uint16)
		}

		statusCount, ok := onlineCount[presence.Status]
		if !ok {
			statusCount = 0
		}
		statusCount++
		onlineCount[presence.Status] = statusCount
		onlineCounts[c.GuildID] = onlineCount
	}
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
		message := getStatsStr(m.GuildID) + "\n"

		for status, count := range onlineCounts[m.GuildID] {
			message += "\n* **" + string(status) + "**: " + strconv.Itoa(int(count))
		}

		if _, err := s.ChannelMessageSendReply(m.ChannelID, message, m.MessageReference); err != nil {
			log.Println("Error:", err)
		}
	}
}
