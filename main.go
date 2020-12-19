package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	configo "github.com/jxsl13/simple-configo"
)

var (
	cfg = &Config{}

	// context stuff
	globalCtx, globalCancel = context.WithCancel(context.Background())
)

func init() {
	var env map[string]string
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}

	if err := configo.Parse(cfg, env); err != nil {
		log.Fatal(err)
	}
}

func main() {
	responses := make(chan ResponseDTO, 128)
	go restConsumer(globalCtx, cfg, responses)
	go discordAnnouncer(globalCtx, cfg, responses)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	globalCancel()

	log.Println("Shutting down, please wait...")
}
