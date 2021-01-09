package service

import (
	"github.com/jxsl13/Teeworlds-Clan-Activity/service/registry"
)

var (
	// new services that want to consume from the playerlist need to registers a new channel here,
	// then they can fetch the date from their channel after the specified period of time that is waited before
	// a new api call is made and the results are fd to the channels.
	playerListConsumerRegistry = registry.NewPlayerListConsumerRegistry()
)
