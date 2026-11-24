package catalog

import (
	"gamedex/internal/api"
	"slices"
	"strings"
)

func Filter(games []api.Game, genre, query string) []api.Game {
	result := []api.Game{}
	for _, g := range games {
		if genre != "" && !slices.Contains(g.Genres, genre) {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(g.Name), strings.ToLower(query)) {
			continue
		}
		result = append(result, g)
	}
	return result
}

func FindGame(games []api.Game, id int) (api.Game, bool) {
	for _, g := range games {
		if g.ID == id {
			return g, true
		}
	}
	return api.Game{}, false
}

func GroupTournamentsByGame(tournaments []api.Tournament) map[int][]api.Tournament {
	byGame := make(map[int][]api.Tournament)
	for _, t := range tournaments {
		byGame[t.GameID] = append(byGame[t.GameID], t)
	}
	return byGame
}