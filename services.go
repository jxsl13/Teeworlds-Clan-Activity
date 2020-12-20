package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	log.Printf("Fretching data every %s\n", cfg.RefreshInterval.String())
	ticker := time.NewTicker(cfg.RefreshInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down REST API consumer...")
		case <-ticker.C:
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
	myMessageID := msg.ID

	err = cleanupBefore(dg, channelID, myMessageID)
	if err != nil {
		log.Printf("Failed to cleanup messages before the initial message: %v", err)
	}

	for {
		select {
		case responseDTO := <-playerChannel:

			err = cleanupAfter(dg, channelID, myMessageID)
			if err != nil {
				log.Printf("Failed to cleanup newer messages: %v", err)
			}

			var sb strings.Builder
			sb.Grow(len(responseDTO.Players) * 64)
			sb.WriteString(fmt.Sprintf("Date: %s\n", time.Now().Format("02.01.2006 15:04:05")))
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
				log.Printf("Error while editing message: %v\n", err)
				continue
			}
		case <-ctx.Done():
			log.Println("Shutting down the posting service...")
			return
		}
	}
}
