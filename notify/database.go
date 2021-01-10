package notify

import (
	"sort"
	"sync"

	"github.com/jxsl13/Teeworlds-Clan-Activity/dto"
)

// Notification contains a player and the discord requestors
// that wanted to be notified
type Notification struct {
	Player     dto.Player
	Requestors []string
}

// Notifications is a list of notifications
type Notifications []Notification

type byNotificationNickname Notifications

func (a byNotificationNickname) Len() int           { return len(a) }
func (a byNotificationNickname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byNotificationNickname) Less(i, j int) bool { return a[i].Player.Name < a[j].Player.Name }

type database struct {
	// ingame nick -> requesting discord user
	requests map[string]map[string]bool

	// discord user -> requested nicks
	requestingUsers map[string]map[string]bool
	mu              sync.RWMutex
}

func (db *database) RequestNotification(discordUser, ingameNick string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.requests[ingameNick] == nil {
		db.requests[ingameNick] = make(map[string]bool)
	}

	if db.requestingUsers[discordUser] == nil {
		db.requestingUsers[discordUser] = make(map[string]bool)
	}

	db.requests[ingameNick][discordUser] = true
	db.requestingUsers[discordUser][ingameNick] = true
}

func (db *database) GetNotificationRequests(discordUser string) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return setToSortedList(db.requestingUsers[discordUser])
}

func (db *database) deleteNotificationRequest(discordUser, ingameNick string) {
	delete(db.requests[ingameNick], discordUser)
	if len(db.requests[ingameNick]) == 0 {
		delete(db.requests, ingameNick)
	}

	delete(db.requestingUsers[discordUser], ingameNick)
	if len(db.requestingUsers[discordUser]) == 0 {
		delete(db.requestingUsers, discordUser)
	}
}

func (db *database) DeleteNotificationRequest(discordUser, ingameNick string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.deleteNotificationRequest(discordUser, ingameNick)
}

func (db *database) ClearNotificationRequests(discordUser string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// remove user from nicknames
	requestedNicks := db.requestingUsers[discordUser]
	for nick := range requestedNicks {
		delete(db.requests[nick], discordUser)

		if len(db.requests[nick]) == 0 {
			delete(db.requests, nick)
		}
	}

	// remove user and his nickname map
	delete(db.requestingUsers, discordUser)
}

func (db *database) isRequestedNick(nick string) bool {
	requestors, ok := db.requests[nick]
	if !ok {
		return false
	}
	if len(requestors) == 0 {
		return false
	}
	return true
}

func (db *database) isRequested(player dto.Player) bool {
	return db.isRequestedNick(player.Name)
}

func (db *database) getRequestors(player dto.Player) []string {
	return setToSortedList(db.requests[player.Name])
}

func (db *database) ParseResponse(response dto.Response) Notifications {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make(Notifications, 0)

	for _, player := range response.Players {
		if db.isRequested(player) {
			requestors := db.getRequestors(player)
			result = append(result, Notification{
				Player:     player,
				Requestors: requestors,
			})

			// after notifying we do not want the user to get any more notifications
			for _, requestor := range requestors {
				db.deleteNotificationRequest(requestor, player.Name)
			}
		}

	}
	sort.Sort(byNotificationNickname(result))
	return result
}
