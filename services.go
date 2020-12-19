package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func restConsumer(ctx context.Context, cfg *Config, rChan chan<- ResponseDTO) {
	client := &http.Client{}

	target, err := url.Parse("https://api.status.tw/2.0/player/list/")
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	// Query params
	params := url.Values{}
	params.Add("clan", cfg.Clanname)
	target.RawQuery = params.Encode()

	log.Println("Request API: ", target.String())

	req, err := http.NewRequest("GET", target.String(), nil)
	if err != nil {
		log.Fatalf("Malformed request, shutting down....")
	}

	ticker := time.NewTicker(cfg.RefreshInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down REST API consumer...")
		case <-ticker.C:
			log.Println("Fetching data...")
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error while fetching data from the status.tw REST API: %v\n", err)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Failed to read response body.")
				log.Printf("BODY: %s\n", string(body))
				resp.Body.Close()
				continue
			}
			responseDTO := ResponseDTO{}
			err = json.Unmarshal(body, &responseDTO)
			if err != nil {
				log.Println("Failed to unmarshal json body response.")
				log.Printf("BODY: %s\n", string(body))
				resp.Body.Close()
				continue
			}
			rChan <- responseDTO
		}
	}
}

func discordAnnouncer(ctx context.Context, cfg *Config, playerChannel <-chan ResponseDTO) {

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
	messages, err := dg.ChannelMessages(channelID, 100, msg.ID, "", "")
	if err != nil {
		log.Fatalf("Failed to fetch messages: %v\n", err)
	}
	if len(messages) > 0 {
		ids := make([]string, 0, len(messages))
		for _, m := range messages {
			ids = append(ids, m.ID)
		}
		err = dg.ChannelMessagesBulkDelete(channelID, ids)
		if err != nil {
			log.Fatalf("Could not delete old messages: %v\n", err)
		}
	}

	myMessageID := msg.ID
	for {
		select {
		case responseDTO := <-playerChannel:
			var sb strings.Builder
			sb.Grow(len(responseDTO.Players) * 50)
			for _, p := range responseDTO.Players {
				sb.WriteString(p.String())
			}

			_, err = dg.ChannelMessageEdit(channelID, myMessageID, sb.String())
			if err != nil {
				log.Printf("Error while editing message: %v\n", err)
				continue
			}
		case <-ctx.Done():
			log.Println("Shutting down the posting service...")
			return
		}
	}
}
