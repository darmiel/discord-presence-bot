package internal

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

func GuildMembersChunk(_ *discordgo.Session, c *discordgo.GuildMembersChunk) {
	for _, presence := range c.Presences {
		onlineCount, ok := OnlineCounts[c.GuildID]
		if !ok {
			onlineCount = make(map[discordgo.Status]uint16)
		}

		statusCount, ok := onlineCount[presence.Status]
		if !ok {
			statusCount = 0
		}
		statusCount++
		onlineCount[presence.Status] = statusCount
		OnlineCounts[c.GuildID] = onlineCount
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore self
	if m.Author.ID == s.State.User.ID {
		log.Println("[Debug] Ignored message from bot")
		return
	}

	if m.Content == "!crawl" {
		_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" | Aight! Crawling ...")
		if err != nil {
			log.Println("Error:", err)
			return
		}
		CheckMembers(s)
	} else if m.Content == "!online" {
		message := getStatsStr(m.GuildID) + "\n"

		for status, count := range OnlineCounts[m.GuildID] {
			message += "\n* **" + string(status) + "**: " + strconv.Itoa(int(count))
		}

		if _, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" | "+message); err != nil {
			log.Println("Error:", err)
		}
	}
}
