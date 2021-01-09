package dto

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jxsl13/Teeworlds-Clan-Activity/markdown"
)

// Player represents a player object
type Player struct {
	Name      string      `json:"name"`
	Clan      string      `json:"clan"`
	Country   int         `json:"country"`
	Score     int         `json:"score"`
	Type      int         `json:"type"`
	FirstSeen string      `json:"first_seen"`
	LastSeen  string      `json:"last_seen"`
	Server    *ServerData `json:"server_data"`
}

func (p *Player) String() string {
	clanFmtStr := fmt.Sprintf("%%-%ds", len(p.Clan))
	name := markdown.WrapInInlineCodeBlock(fmt.Sprintf("%-16s", p.Name))
	clan := markdown.WrapInInlineCodeBlock(fmt.Sprintf(clanFmtStr, p.Clan))

	servername := ""
	if p.Server != nil {
		servername = p.Server.Name
	}

	servername = markdown.WrapInInlineCodeBlock(servername)

	return fmt.Sprintf("%s %s %s on %s\n", markdown.Flag(p.Country), name, clan, servername)
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
	longestServerName := 0
	for _, player := range p {
		if len([]rune(player.Name)) > longestName {
			longestName = len([]rune(player.Name))
		}
		if player.Server != nil && len([]rune(player.Server.Name)) > longestServerName {
			longestServerName = len([]rune(player.Server.Name))
		}
	}

	sfl = make([]PlayerStringTuple, 0, len(p))
	for _, player := range p {
		nameFmtStr := fmt.Sprintf("%%-%ds", longestName)

		name := markdown.WrapInInlineCodeBlock(fmt.Sprintf(nameFmtStr, player.Name))

		servername := "(unknown)"

		if player.Server != nil {
			servername = player.Server.Name
		}

		// alignment
		serverFmtStr := fmt.Sprintf("%%-%ds", longestServerName)
		// wrap in monospaced inline code block
		servername = markdown.WrapInInlineCodeBlock(fmt.Sprintf(serverFmtStr, servername))

		sfl = append(sfl, PlayerStringTuple{
			Player: player,
			String: fmt.Sprintf("%s %s on %s\n", markdown.Flag(player.Country), name, servername),
		})

	}

	sort.Sort(ByPlayerName(sfl))
	return sfl
}
