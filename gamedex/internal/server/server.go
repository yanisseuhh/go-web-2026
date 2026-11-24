package server

import (
	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	Games       []api.Game
	Tournaments map[int][]api.Tournament
	tmpl        *template.Template
}

type ListePage struct {
	Games []api.Game
	Query string
}

type JeuPage struct {
	Game        api.Game
	Tournaments []api.Tournament
}

type ErreurPage struct {
	Code    int
	Message string
}

func New(games []api.Game, tournaments []api.Tournament) *Server {
	return &Server{
		Games:       games,
		Tournaments: catalog.GroupTournamentsByGame(tournaments),
		tmpl:        template.Must(template.ParseGlob("web/templates/*.html")),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/jeux", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /jeux", s.list)
	mux.HandleFunc("GET /jeux/{id}", s.detail)

	return mux
}

func (s *Server) show(w http.ResponseWriter, page string, data any) {
	if err := s.tmpl.ExecuteTemplate(w, page, data); err != nil {
		log.Println(err)
	}
}

func (s *Server) renderError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	s.show(w, "erreur.html", ErreurPage{
		Code:    code,
		Message: message,
	})
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	games := catalog.Filter(s.Games, "", q)

	s.show(w, "liste.html", ListePage{
		Games: games,
		Query: q,
	})
}

func (s *Server) detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		s.renderError(w, http.StatusBadRequest, "L'identifiant du jeu doit être un nombre positif.")
		return
	}

	game, ok := catalog.FindGame(s.Games, id)
	if !ok {
		s.renderError(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	s.show(w, "jeu.html", JeuPage{
		Game:        game,
		Tournaments: s.Tournaments[id],
	})
}