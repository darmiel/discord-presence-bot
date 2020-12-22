package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func CheckMembers(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		log.Println("+ Guild:", guild.ID)

		// Get online count
		if count, ok := OnlineCounts[guild.ID]; ok {
			tm, ok := CrawlStarts[guild.ID]
			if !ok {
				tm = 0
			}

			save(guild.ID, tm, count)
		}

		// reset count
		OnlineCounts[guild.ID] = make(map[discordgo.Status]uint16)
		CrawlStarts[guild.ID] = time.Now().Unix()

		if err := s.RequestGuildMembers(Guild, "", 200, true); err != nil {
			log.Println("[WARN] Error requesting guild members:", err)
			return
		}
	}
}

func cumulativeCounts(guildId string) (res uint) {
	onlineCount, ok := OnlineCounts[guildId]
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

	crawlStart, ok := CrawlStarts[guildId]
	if !ok {
		crawlStart = unixNow
	}

	return uint32(unixNow - crawlStart), cumulativeCounts(guildId)
}

func getStatsStr(guildId string) (res string) {
	timespan, online := getStats(guildId)
	return fmt.Sprintf("Online between the last %v seconds: %v", timespan, online)
}
