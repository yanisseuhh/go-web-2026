package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", 8080, "port du serveur")
	flag.Parse()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Serveur sur http://localhost%s\n", addr)
	http.ListenAndServe(addr, mux)
}

func accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `<!DOCTYPE html>
<html lang="fr">
<head>
	<meta charset="utf-8">
	<title>GameDex</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<h1>Bienvenue sur GameDex</h1>
	<p><a href="/jeux">Voir la liste des jeux</a></p>
</body>
</html>`)
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux")
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "identifiant invalide", http.StatusBadRequest)
		return
	}
	if id > 12 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Fiche du jeu %d\n", id)
}