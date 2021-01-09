package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jxsl13/Teeworlds-Clan-Activity/config"
	"github.com/jxsl13/Teeworlds-Clan-Activity/markdown"
	"github.com/jxsl13/Teeworlds-Clan-Activity/notify"
)

// OnlinePlayersNotifier announces to the configured discord channed an updated message that
// contains a list of the currently online members
func OnlinePlayersNotifier(ctx context.Context, cfg *config.Config) {
	const KEY = "OnlinePlayersNotifier"
	playerListConsumerRegistry.Register(KEY, 64)
	defer playerListConsumerRegistry.Unregister(KEY)

	channelID := cfg.NotificationChannel

	dg, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		log.Println(err)
		return
	}

	// handle '?notify <nickname>'
	dg.AddHandler(onlinePlayerNotificationRequest(channelID))

	err = dg.Open()
	if err != nil {
		log.Fatalln("error: could not establish a connection to the discord api, please check your credentials")
		return
	}
	defer dg.Close()

	msg, err := dg.ChannelMessageSend(channelID, "Write `?notify <nickname>` in order to get a notification when that player is online.")
	if err != nil {
		log.Printf("Failed to fetch the configured channelID, please try again: %s", err)
		return
	}
	myMessageID := msg.ID

	err = cleanupBefore(dg, channelID, myMessageID)
	if err != nil {
		log.Printf("Failed to cleanup messages before the initial message: %v", err)
	}

	for {
		select {
		case responseDTO := <-playerListConsumerRegistry.Get(KEY):
			notifications := notify.ParseResponse(responseDTO)

			var sb strings.Builder
			sb.Grow(len(notifications) * 64)

			for _, n := range notifications {
				player := markdown.WrapInInlineCodeBlock(n.Player.Name)
				server := markdown.WrapInInlineCodeBlock(n.Player.Server.Name)
				header := fmt.Sprintf("%s joined the server %s\n", player, server)
				sb.WriteString(header)

				for idx, requestor := range n.Requestors {
					if idx < len(n.Requestors) {
						sb.WriteString(fmt.Sprintf("%s, ", mentionUser(requestor)))
					} else {
						sb.WriteString(fmt.Sprintf("%s", mentionUser(requestor)))
					}
				}

				_, err = dg.ChannelMessageSend(channelID, sb.String())
				if err != nil {
					log.Printf("Failed to send notification message: %s", err)
					continue
				}
			}
		case <-ctx.Done():
			log.Println("Shutting down the posting service...")
			return
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func onlinePlayerNotificationRequest(notificationChannelID string) func(s *discordgo.Session, m *discordgo.MessageCreate) {

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID != notificationChannelID {
			return
		}

		text := m.Content
		if len(text) <= 1 {
			return
		}

		if text[0] != '?' {
			return
		}
		author := m.Author.ID
		text = text[1:]
		tokens := strings.SplitN(text, " ", 2)
		if len(tokens) < 2 {
			return
		}
		command := tokens[0]
		nickname := tokens[1]

		if command != "notify" {
			return
		}

		log.Printf("User: %s requested to be notified when '%s' is online.", m.Author.String(), nickname)
		notify.RequestNotification(author, nickname)
		msgText := fmt.Sprintf("%s, you will be notified once '%s' joins a server.", m.Author.Mention(), nickname)
		if _, err := s.ChannelMessageSend(notificationChannelID, msgText); err != nil {
			log.Printf("Failed to send notification request received confirmation.")
		}
	}
}
