package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// ResponseDTO is a data transfer object that maps
type ResponseDTO struct {
	Players Players `json:"players"`
}

// ServerData that is attached to a player
type ServerData struct {
	IP         string `json:"server_ip"`
	Port       string `json:"server_port"`
	FirstSeen  string `json:"first_seen"`
	LastSeen   string `json:"last_seen"`
	Version    string `json:"version"`
	Name       string `json:"name"`
	Password   bool   `json:"password"`
	SkillLevel int    `json:"server_level"`
	NumPlayers int    `json:"num_players"`
	MaxPlayers int    `json:"max_players"`
	GameType   string `json:"gamemode"`
	Map        string `json:"map"`
	Country    string `json:"country"`
	Master     string `json:"master"`
}

// Player represents a player object
type Player struct {
	Name      string     `json:"name"`
	Clan      string     `json:"clan"`
	Country   int        `json:"country"`
	Score     int        `json:"score"`
	Type      int        `json:"type"`
	FirstSeen string     `json:"first_seen"`
	LastSeen  string     `json:"last_seen"`
	Server    ServerData `json:"server_data"`
}

func (p *Player) String() string {
	clanFmtStr := fmt.Sprintf("%%-%ds", len(p.Clan))
	name := WrapInInlineCodeBlock(fmt.Sprintf("%-16s", p.Name))
	clan := WrapInInlineCodeBlock(fmt.Sprintf(clanFmtStr, p.Clan))

	servername := p.Server.Name
	l := 36
	if len(servername) < 36 {
		l = len(servername)
	}
	servername = WrapInInlineCodeBlock(servername[:l])

	return fmt.Sprintf("%s %s %s on %s\n", Flag(p.Country), name, clan, servername)
}

// LastSeenIn the last x minutes, seconds, hours
func (p *Player) LastSeenIn(d time.Duration) bool {
	t, err := time.Parse("2006-01-02 15:04:05", p.LastSeen)
	if err != nil {
		log.Printf("Malformed time string LastSeen: %s", p.LastSeen)
		return true
	}
	return time.Since(t) < d
}

// Players is a list of players
type Players []Player

// PlayerStringTuple Can be sorted
type PlayerStringTuple struct {
	Player Player
	String string
}

// ByPlayerName is used to sort the result
type ByPlayerName []PlayerStringTuple

func (a ByPlayerName) Len() int      { return len(a) }
func (a ByPlayerName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPlayerName) Less(i, j int) bool {
	return a[i].Player.Name < a[j].Player.Name
}

// StringFormatList returns a list of properly formated strings based on the length of all
// Player's name , clan, etc lengths.
func (p Players) StringFormatList() (sfl []PlayerStringTuple) {
	var sb strings.Builder
	sb.WriteString("")

	longestName := 0
	longestClan := 0
	longestServerName := 0
	for _, player := range p {
		if len(player.Name) > longestName {
			longestName = len(player.Name)
		}
		if len(player.Clan) > longestClan {
			longestClan = len(player.Clan)
		}
		if len(player.Server.Name) > longestServerName {
			longestServerName = len(player.Server.Name)
		}
	}

	sfl = make([]PlayerStringTuple, 0, len(p))
	for _, player := range p {
		nameFmtStr := fmt.Sprintf("%%-%ds", longestName)
		clanFmtStr := fmt.Sprintf("%%-%ds", longestClan)

		name := WrapInInlineCodeBlock(fmt.Sprintf(nameFmtStr, player.Name))
		clan := WrapInInlineCodeBlock(fmt.Sprintf(clanFmtStr, player.Clan))

		if longestServerName > 36 {
			longestServerName = 36
		}
		servername := player.Server.Name
		serverFmtStr := fmt.Sprintf("%%-%ds", longestServerName)
		servername = WrapInInlineCodeBlock(fmt.Sprintf(serverFmtStr, servername[:longestServerName]))

		sfl = append(sfl, PlayerStringTuple{
			Player: player,
			String: fmt.Sprintf("%s %s %s on %s\n", Flag(player.Country), name, clan, servername),
		})

	}

	sort.Sort(ByPlayerName(sfl))
	return sfl
}
