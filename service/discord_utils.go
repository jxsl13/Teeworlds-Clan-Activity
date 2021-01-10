package service

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func mentionUser(ID string) string {
	return "<@" + ID + ">"
}

// SplitChannelMessageSend properly splits long output in order to accepted by the discord servers.
// also properly wrap single codeblocks that were split during this process
func SplitChannelMessageSend(s *discordgo.Session, channelID, text string) {
	const codeblockDelimiter = "```"

	codeblockFound := strings.Count(text, codeblockDelimiter) == 2
	chunks := Split(text, "\n", 1800)

	codeblockBegin := -1
	codeblockEnd := -1

	if codeblockFound {
		beginSet := false

		for idx, chunk := range chunks {

			delimiterCount := strings.Count(chunk, codeblockDelimiter)
			if delimiterCount == 2 {
				codeblockBegin = idx
				codeblockEnd = idx
				break
			} else if delimiterCount == 1 && !beginSet {
				codeblockBegin = idx
				beginSet = true
			} else if delimiterCount == 1 && beginSet {
				codeblockEnd = idx
				break
			}

		}
	}

	isCodeblockSplit := 0 <= codeblockBegin && codeblockBegin < codeblockEnd

	for idx, chunk := range chunks {

		if isCodeblockSplit {
			if idx == codeblockBegin {
				chunk = chunk + codeblockDelimiter
			} else if codeblockBegin < idx && idx < codeblockEnd {
				chunk = codeblockDelimiter + chunk + codeblockDelimiter
			} else if idx == codeblockEnd {
				chunk = codeblockDelimiter + chunk
			}
		}

		if _, err := s.ChannelMessageSend(channelID, chunk); err != nil {
			log.Println(err)
		}
	}
}

func cleanupMsgs(dg *discordgo.Session, channelID string, limit int, beforeMsgID, afterMsgID, aroundMsgID string) error {
	messages, err := dg.ChannelMessages(channelID, limit, beforeMsgID, afterMsgID, aroundMsgID)
	if err != nil {
		log.Printf("Failed to fetch messages: %v\n", err)
	}

	for len(messages) > 0 {
		ids := make([]string, 0, len(messages))
		for _, m := range messages {
			ids = append(ids, m.ID)
		}
		err = dg.ChannelMessagesBulkDelete(channelID, ids)
		if err != nil {
			log.Printf("Could not bulk delete messages, trying individual deletion.")
			for _, msgID := range ids {
				err = dg.ChannelMessageDelete(channelID, msgID)
				if err != nil {
					return fmt.Errorf("Unable to cleanup messages individually: %w", err)
				}
			}
		}

		messages, err = dg.ChannelMessages(channelID, limit, beforeMsgID, afterMsgID, aroundMsgID)
		if err != nil {
			return fmt.Errorf("Failed to fetch messages: %v", err)
		}
	}

	return nil
}

func cleanupBefore(dg *discordgo.Session, channelID, messageID string) error {
	return cleanupMsgs(dg, channelID, 100, messageID, "", "")
}

func cleanupAfter(dg *discordgo.Session, channelID, messageID string) error {
	return cleanupMsgs(dg, channelID, 100, "", messageID, "")
}

func sleepAndDelete(dg *discordgo.Session, channelID, messageID string, duration time.Duration) {
	time.Sleep(duration)
	dg.ChannelMessageDelete(channelID, messageID)
}
