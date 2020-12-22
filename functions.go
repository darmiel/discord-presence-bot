package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func checkMembers(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		log.Println("+ Guild:", guild.ID)

		// Get online count
		if count, ok := onlineCounts[guild.ID]; ok {
			tm, ok := crawlStarts[guild.ID]
			if !ok {
				tm = 0
			}

			save(guild.ID, tm, count)
		}

		// reset count
		onlineCounts[guild.ID] = make(map[discordgo.Status]uint16)
		crawlStarts[guild.ID] = time.Now().Unix()

		if err := s.RequestGuildMembers(Guild, "", 200, true); err != nil {
			log.Println("[WARN] Error requesting guild members:", err)
			return
		}
	}
}

func cumulativeCounts(guildId string) (res uint) {
	onlineCount, ok := onlineCounts[guildId]
	var count uint = 0
	if ok {
		for _, c := range onlineCount {
			count += uint(c)
		}
	}
	return count
}

func getStats(guildId string) (timespan uint32, online uint) {
	unixNow := time.Now().Unix()

	crawlStart, ok := crawlStarts[guildId]
	if !ok {
		crawlStart = unixNow
	}

	return uint32(unixNow - crawlStart), cumulativeCounts(guildId)
}

func getStatsStr(guildId string) (res string) {
	timespan, online := getStats(guildId)
	return fmt.Sprintf("Online between the last %v seconds: %v", timespan, online)
}
