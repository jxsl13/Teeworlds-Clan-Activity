package main

import (
	"fmt"
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
	return fmt.Sprintf("```%16s %12s Last seen: %v```", p.Name, p.Clan, p.LastSeen)
}
