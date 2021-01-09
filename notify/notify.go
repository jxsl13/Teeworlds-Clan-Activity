package notify

import "github.com/jxsl13/Teeworlds-Clan-Activity/dto"

var db database = database{
	// using 64, as that's currently the max number
	// of ingame players and maybe also the max number
	// of players in a clan
	requests:        make(map[string]map[string]bool, 64),
	requestingUsers: make(map[string]map[string]bool, 64),
}

// RequestNotification allows a discord user to request to be
// notified when a specific user joins the game, this does only handle the notification logic
func RequestNotification(DiscordUser, Nickname string) {
	db.RequestNotification(DiscordUser, Nickname)
}

// GetNotificationRequests allow you to fetch all requested notifications of a single discord user.
func GetNotificationRequests(DiscordUser string) []string {
	return db.GetNotificationRequests(DiscordUser)
}

// DeleteNotificationRequest allows a discord user to delete a specific notification request
// by specifying a nickname
func DeleteNotificationRequest(DiscordUser, Nickname string) {
	db.DeleteNotificationRequest(DiscordUser, Nickname)
}

// ClearNotificationRequests allows to clear all notification requests of a specific discord user.
func ClearNotificationRequests(DiscordUser string) {
	db.ClearNotificationRequests(DiscordUser)
}

// ParseResponse parses the response dto object and creates per
// ingame nick a Notification object that contains all of the
// player data as well as a list of people that wanted to be notified
// when that player joined the game
func ParseResponse(Response dto.Response) Notifications {
	return db.ParseResponse(Response)
}
