package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
	txt := "This is a one tme notification channel. Once the requested player joins the game, you will not receive any more notifications after your initial one.\nWrite `?notify <nickname>` in order to get a notification when that player is online.\nUse `?list` to see all of your notification requests.\nAnd execute `?unnotify_all` to delete all of your notification requests."

	msg, err := dg.ChannelMessageSend(channelID, txt)
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
				header := ""
				if n.Player.Server == nil {
					header = fmt.Sprintf("%s joined the server %s\n", player, server)
				} else {
					header = fmt.Sprintf("%s joined the game.\n", player)
				}

				sb.WriteString(header)

				for idx, requestor := range n.Requestors {
					if idx < len(n.Requestors)-1 {
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
		if m.Author.ID == s.State.User.ID {
			return
		}
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
		if len(tokens) < 1 {
			return
		}
		command := tokens[0]

		args := ""
		if len(tokens) > 1 {
			args = tokens[1]
		}

		switch command {
		case "notify":
			nickname := args

			notify.RequestNotification(author, nickname)

			// delete request message
			err := s.ChannelMessageDelete(notificationChannelID, m.ID)
			if err != nil {
				log.Printf("Failed to delete request message of 'notify': %s", err)
			}
			// logging
			log.Printf("'%s' requested to be notified when '%s' is online.", m.Author.String(), nickname)

			// discord answer
			msgText := fmt.Sprintf("%s, you will be notified once %s joins a server.", m.Author.Mention(), markdown.WrapInInlineCodeBlock(nickname))
			msg, err := s.ChannelMessageSend(notificationChannelID, msgText)
			if err != nil {
				log.Printf("Error when sending confirmation message to notification channel: %s", err)
			}

			sleepAndDelete(s, notificationChannelID, msg.ID, 5*time.Second)
		case "unnotify_all":
			notify.ClearNotificationRequests(author)

			// delete request message
			err := s.ChannelMessageDelete(notificationChannelID, m.ID)
			if err != nil {
				log.Printf("Failed to delete request message of 'unnotify_all': %s", err)
			}

			log.Printf("'%s' requested to delete all notification requests.", m.Author.String())
			txt := fmt.Sprintf("%s, your notification requests were deleted.", m.Author.Mention())
			msg, err := s.ChannelMessageSend(notificationChannelID, txt)
			if err != nil {
				log.Printf("Failed to delete notify deletion request message: %s", err)
			}
			sleepAndDelete(s, notificationChannelID, msg.ID, 5*time.Second)
		case "list":
			list := notify.GetNotificationRequests(author)
			header := fmt.Sprintf("%s, your notification requests are:\n", m.Author.Mention())

			var sb strings.Builder
			sb.Grow(len(header) + 18*len(list))

			sb.WriteString(header)

			for _, nickname := range list {
				sb.WriteString(markdown.WrapInInlineCodeBlock(nickname))
				sb.WriteString("\n")
			}

			msg, err := s.ChannelMessageSend(notificationChannelID, sb.String())
			if err != nil {
				log.Printf("Failed to delete notify deletion request message: %s", err)
			}
			sleepAndDelete(s, notificationChannelID, msg.ID, 15*time.Second)
		default:
			return
		}

	}
}
