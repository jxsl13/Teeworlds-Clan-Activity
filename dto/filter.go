package dto

// FilterByClan returns a new Response that only contains Player swith a matching clan
func FilterByClan(clan string, in Response) Response {
	clanMembers := make(Players, 0, 1)
	for _, player := range in.Players {
		if player.Clan == clan {
			clanMembers = append(clanMembers, player)
		}
	}
	return Response{clanMembers}
}
