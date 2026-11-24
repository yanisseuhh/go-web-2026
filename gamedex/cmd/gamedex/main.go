package main

import (
	"flag"
	"fmt"
	"gamedex/internal/api"
	"gamedex/internal/server"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "port du serveur")
	flag.Parse()

	games, err := api.FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux : ", err)
	}

	tournaments, err := api.FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois : ", err)
	}

	srv := server.New(games, tournaments)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("GameDex sur http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}