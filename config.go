package main

import (
	"time"

	configo "github.com/jxsl13/simple-configo"
)

// Config is the discord bot configuration representation
type Config struct {
	DiscordBotToken string
	DiscordChannel  string
	Owner           string

	Clanname        string
	Gametype        string
	RefreshInterval time.Duration
}

// Name is the name of the configuration
func (c *Config) Name() string {
	return "Teeworlds Clan Activity"
}

// Options returns a list of for this configuration required options that can be parsed by
// simple-configo
func (c *Config) Options() configo.Options {
	optionsList := configo.Options{
		{
			Key:           "TCA_DISCORD_BOT_TOKEN",
			Mandatory:     true,
			Description:   "Token that can be retrieved by creating an app here: https://discord.com/developers/applications",
			ParseFunction: configo.DefaultParserString(&c.DiscordBotToken),
		},
		{
			Key:           "TCA_DISCORD_CHANNEL",
			Mandatory:     true,
			Description:   "After you gave the bot access to your discord server via the discord website, you cna get a channel ID from your server that can be used here.",
			ParseFunction: configo.DefaultParserString(&c.DiscordChannel),
		},
		{
			Key:           "TCA_TEEWORLDS_CLAN",
			Mandatory:     true,
			Description:   "The clan that should be watched.",
			ParseFunction: configo.DefaultParserString(&c.Clanname),
		},
		{
			Key:           "TCA_TEEWORLDS_GAMETYPE",
			Description:   "The gametype that should be watched.",
			DefaultValue:  "",
			ParseFunction: configo.DefaultParserString(&c.Gametype),
		},
		{
			Key:           "TCA_OWNER",
			Description:   "The Discord user that owns the bot.",
			DefaultValue:  "",
			ParseFunction: configo.DefaultParserRegex(&c.Owner, `.+#\d{4,4}|`, "Require nick#1234 Discord nick."),
		},
		{
			Key:           "TCA_REFRESH_INTERVAL",
			Description:   "At what frequency to request data from the status.tw REST API",
			DefaultValue:  "1m",
			ParseFunction: configo.DefaultParserDuration(&c.RefreshInterval),
		},
	}
	return optionsList
}
