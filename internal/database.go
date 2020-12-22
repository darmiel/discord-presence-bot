package internal

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

var db *sql.DB

func InitDatabase() {
	// Init database
	database, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		log.Fatalln(err)
		return
	}
	db = database

	statement, err := db.Prepare("create table if not exists " + SqlTable + " (" +
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

func save(guildId string, at int64, counts map[discordgo.Status]uint16) {
	log.Println("Saving infos for guild", guildId, "at", at, "with a count of", counts)

	var columns = "guild, crawl_time"
	var values = "'" + guildId + "', " + strconv.Itoa(int(at))

	for key, val := range counts {
		columns += ", `" + string(key) + "`"
		values += ", " + strconv.Itoa(int(val))
	}

	var query = "INSERT INTO `" + SqlTable + "` (" + columns + ") VALUES (" + values + ");"

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
