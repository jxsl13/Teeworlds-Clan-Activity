package dto

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
