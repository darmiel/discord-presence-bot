package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/// Discord
var (
	Token string
	Guild string
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

	// Init database
	database, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		log.Fatalln(err)
		return
	}
	db = database

	statement, err := db.Prepare("create table if not exists " + table + " (" +
		"guild text," +
		"crawl_time integer," +
		"online integer DEFAULT 0 NOT NULL," +
		"idle integer DEFAULT 0 NOT NULL," +
		"dnd integer DEFAULT 0 NOT NULL," +
		"invisible integer DEFAULT 0 NOT NULL," +
		"offline integer DEFAULT 0 NOT NULL," +
		"CONSTRAINT active_user_counts_pk PRIMARY KEY (guild, crawl_time)" +
		");")
	if err != nil {
		log.Fatalln("Statement error:", err)
		return
	}

	if _, err := statement.Exec(); err != nil {
		log.Fatalln("Statement error 2:", err)
		return
	}
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

/////////////////////////////////////////////////////////////////////////////////////

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

func save(guildId string, at int64, counts map[discordgo.Status]uint16) {
	log.Println("Saving infos for guild", guildId, "at", at, "with a count of", counts)

	var columns = "guild, crawl_time"
	var values = "'" + guildId + "', " + strconv.Itoa(int(at))

	for key, val := range counts {
		columns += ", `" + string(key) + "`"
		values += ", " + strconv.Itoa(int(val))
	}

	var query = "INSERT INTO `" + table + "` (" + columns + ") VALUES (" + values + ");"

	statement, err := db.Prepare(query)
	if err != nil {
		log.Println("[WARN] Statement error:", err)
		return
	}

	if _, err := statement.Exec(); err != nil {
		log.Println("[WARN] Exec error:", err)
	}

	log.Println("[DATABASE] Saved infos! Query:", query)
}

/////////////////////////////////////////////////////////////////////////////////////

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
