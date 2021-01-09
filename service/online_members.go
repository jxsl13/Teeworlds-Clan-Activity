package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jxsl13/Teeworlds-Clan-Activity/config"
	"github.com/jxsl13/Teeworlds-Clan-Activity/dto"
	"github.com/jxsl13/Teeworlds-Clan-Activity/markdown"
)

// OnlineMembersAnnouncer announces to the configured discord channed an updated message that
// contains a list of the currently online members
func OnlineMembersAnnouncer(ctx context.Context, cfg *config.Config) {
	const KEY = "OnlineMembersAnnouncer"
	playerListConsumerRegistry.Register(KEY, 64)
	defer playerListConsumerRegistry.Unregister(KEY)

	dg, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		log.Println(err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatalln("error: could not establish a connection to the discord api, please check your credentials")
		return
	}
	defer dg.Close()

	channelID := cfg.DiscordChannel
	msg, err := dg.ChannelMessageSend(channelID, "I'm back my dudines and dudettes aaaand my dudes!")
	if err != nil {
		log.Fatalln("Failed to fetch the configured channelID, please try again.")
	}
	myMessageID := msg.ID

	err = cleanupBefore(dg, channelID, myMessageID)
	if err != nil {
		log.Printf("Failed to cleanup messages before the initial message: %v", err)
	}

	retries := 0

	for {
		select {
		case responseDTO := <-playerListConsumerRegistry.Get(KEY):

			err = cleanupAfter(dg, channelID, myMessageID)
			if err != nil {
				log.Printf("Failed to cleanup newer messages: %v", err)
			}

			// only get clan members from the list
			responseDTO = dto.FilterByClan(cfg.Clanname, responseDTO)

			var sb strings.Builder
			sb.Grow(len(responseDTO.Players) * 64)
			header := fmt.Sprintf("Date: %s\nClan: %s\n\n",
				time.Now().Format("02.01.2006 15:04:05"),
				markdown.WrapInInlineCodeBlock(cfg.Clanname),
			)
			sb.WriteString(header)
			for _, t := range responseDTO.Players.StringFormatList() {
				p := t.Player
				s := t.String

				if p.LastSeenIn(cfg.RefreshInterval) {
					sb.WriteString(s)
				}
			}
			sb.WriteString("\n")

			_, err = dg.ChannelMessageEdit(channelID, myMessageID, sb.String())
			if err != nil {

				msg, err = dg.ChannelMessageSend(channelID, sb.String())
				if err != nil {
					retries++
					if retries > 60 {
						log.Printf("Could not create new message after old one was deleted: %s", err.Error())
						log.Println("Reached 60 retries: Shutting down announcer service...")
						return
					}
				}
				myMessageID = msg.ID
				continue
			}
			retries = 0
		case <-ctx.Done():
			log.Println("Shutting down the posting service...")
			return
		}
	}
}
