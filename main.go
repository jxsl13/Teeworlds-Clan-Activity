package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/jxsl13/Teeworlds-Clan-Activity/config"
	"github.com/jxsl13/Teeworlds-Clan-Activity/service"
	configo "github.com/jxsl13/simple-configo"
)

var (
	cfg = &config.Config{}

	// context stuff
	globalCtx, globalCancel = context.WithCancel(context.Background())
)

func init() {
	// read env
	var env map[string]string
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}

	// parse configuration
	if err := configo.Parse(cfg, env); err != nil {
		log.Fatal(err)
	}
}

func main() {

	go service.PlayerListFetcher(globalCtx, cfg)
	go service.OnlineMembersAnnouncer(globalCtx, cfg)

	// optional feature
	if cfg.NotificationChannel != "" {
		go service.OnlinePlayersNotifier(globalCtx, cfg)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	globalCancel()

	log.Println("Shutting down, please wait...")
}
