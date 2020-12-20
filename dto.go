package main

import (
	"fmt"
	"log"
	"time"
)

// ResponseDTO is a data transfer object that maps
type ResponseDTO struct {
	Players []Player `json:"players"`
}

// Player represents a player object
type Player struct {
	Name      string `json:"name"`
	Clan      string `json:"clan"`
	Country   int    `json:"country"`
	Score     int    `json:"score"`
	Type      int    `json:"type"`
	FirstSeen string `json:"first_seen"`
	LastSeen  string `json:"last_seen"`
}

func (p *Player) String() string {
	clanFmtStr := fmt.Sprintf("%%-%ds", len(p.Clan))
	name := WrapInInlineCodeBlock(fmt.Sprintf("%-22s", p.Name))
	clan := WrapInInlineCodeBlock(fmt.Sprintf(clanFmtStr, p.Clan))

	return fmt.Sprintf("%s %s %s\n", Flag(p.Country), name, clan)
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
