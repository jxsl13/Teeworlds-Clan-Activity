package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jxsl13/Teeworlds-Clan-Activity/config"
	"github.com/jxsl13/Teeworlds-Clan-Activity/dto"
)

// PlayerListFetcher fetches data from the status.tw api and pushed it to all channels in the
// channel map passed as cmap parameter
func PlayerListFetcher(ctx context.Context, cfg *config.Config) {

	client := &http.Client{}

	target, err := url.Parse("https://api.status.tw/2.0/player/list/")
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	log.Println("Request API: ", target.String())

	req, err := http.NewRequest("GET", target.String(), nil)
	if err != nil {
		log.Fatalf("Malformed request, shutting down....: %s", err.Error())
	}

	log.Printf("Fretching data every %s\n", cfg.RefreshInterval.String())
	ticker := time.NewTicker(cfg.RefreshInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down REST API consumer...")
			return
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
			responseDTO := dto.Response{}
			err = json.Unmarshal(body, &responseDTO)
			if err != nil {
				log.Println("Failed to unmarshal json body response.")
				log.Printf("BODY: %s\n", string(body))
				resp.Body.Close()
				continue
			}

			// send to all registered services
			for _, channel := range playerListConsumerRegistry.GetAll() {
				channel <- responseDTO
			}
			resp.Body.Close()
		}
	}
}
