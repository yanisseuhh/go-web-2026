package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

type Tournament struct {
	ID     int     `json:"id"`
	GameID int     `json:"gameId"`
	Name   string  `json:"name"`
	City   string  `json:"city"`
	Date   string  `json:"date"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

type ListePage struct {
	Games []Game
	Query string
}

type JeuPage struct {
	Game        Game
	Tournaments []Tournament
}

type ErreurPage struct {
	Code    int
	Message string
}

var (
	client = &http.Client{Timeout: 3 * time.Second}
	tmpl   = template.Must(template.ParseGlob("templates/*.html"))
)

func FetchGames() ([]Game, error) {
	resp, err := client.Get("http://localhost:8081/api/games")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var games []Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	return games, err
}

func FetchGame(id int) (Game, error) {
	url := fmt.Sprintf("http://localhost:8081/api/games/%d", id)
	resp, err := client.Get(url)
	if err != nil {
		return Game{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Game{}, fmt.Errorf("introuvable")
	}
	if resp.StatusCode != http.StatusOK {
		return Game{}, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var game Game
	err = json.NewDecoder(resp.Body).Decode(&game)
	return game, err
}

func FetchGameTournaments(id int) ([]Tournament, error) {
	url := fmt.Sprintf("http://localhost:8081/api/games/%d/tournaments", id)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var tournaments []Tournament
	err = json.NewDecoder(resp.Body).Decode(&tournaments)
	return tournaments, err
}

func renderError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	err := tmpl.ExecuteTemplate(w, "erreur.html", ErreurPage{
		Code:    code,
		Message: message,
	})
	if err != nil {
		log.Println(err)
	}
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	games, err := FetchGames()
	if err != nil {
		renderError(w, http.StatusBadGateway, "L'API externe est indisponible.")
		return
	}

	q := r.URL.Query().Get("q")
	err = tmpl.ExecuteTemplate(w, "liste.html", ListePage{
		Games: games,
		Query: q,
	})
	if err != nil {
		log.Println(err)
	}
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		renderError(w, http.StatusBadRequest, "L'identifiant du jeu doit être un nombre positif.")
		return
	}

	game, err := FetchGame(id)
	if err != nil {
		if err.Error() == "introuvable" {
			renderError(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		} else {
			renderError(w, http.StatusBadGateway, "L'API externe est indisponible.")
		}
		return
	}

	tournaments, err := FetchGameTournaments(id)
	if err != nil {
		renderError(w, http.StatusBadGateway, "Impossible de charger les tournois.")
		return
	}

	err = tmpl.ExecuteTemplate(w, "jeu.html", JeuPage{
		Game:        game,
		Tournaments: tournaments,
	})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := flag.Int("port", 8080, "port du serveur")
	flag.Parse()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/jeux", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Serveur sur http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}