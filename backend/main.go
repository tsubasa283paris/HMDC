package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tsubasa283paris/HMDC/api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

// Server struct
type Server struct {
	router *chi.Mux
}

// Constructor for struct Server
func New() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

// Executed on server initialization
func (s *Server) Init(env string) {
	log.Println("env: ", env)
}

// Executed before some API
func (s *Server) Middleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(time.Second * 60))
}

// Routing structures
func (s *Server) Router() {
	c := api.NewController()
	s.router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(Auth("db connection")) // middleware

		// users API
		apiRouter.Route("/users", func(users chi.Router) {
			users.Get("/", api.Handler(c.GetUsers).ServeHTTP)
			users.Get("/{userId}/stats", api.Handler(c.GetUserStats).ServeHTTP)
			users.Get("/{userId}/duelhistory", api.Handler(c.GetUserDuelHistory).ServeHTTP)
			users.Get("/{userId}/decks", api.Handler(c.GetUserDecks).ServeHTTP)
			users.Get("/{userId}/details", api.Handler(c.GetUserDetails).ServeHTTP)
			users.Put("/{userId}/details", api.Handler(c.PutUserDetails).ServeHTTP)
		})

		// decks API
		apiRouter.Route("/decks", func(users chi.Router) {
			users.Get("/", api.Handler(c.GetDecks).ServeHTTP)
		})
	})

	// auth API
	s.router.Route("/api/auth", func(authRouter chi.Router) {
		authRouter.Post("/login", api.Handler(c.Login).ServeHTTP)
		authRouter.Post("/signup", api.Handler(c.SignUp).ServeHTTP)
	})
}

// Authentication
func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization") // TODO: implement token system
			// if token != "admin" {
			// 	api.RespondJSON(
			// 		w,
			// 		http.StatusUnauthorized,
			// 		api.ErrorBody{
			// 			Error: "invalid token",
			// 		},
			// 	)
			// 	return
			// }
			userID := token // TODO: acquire user id from token
			r.Header.Set("UserID", userID)
			h.ServeHTTP(w, r)
		})
	}
	return
}

func main() {
	var (
		port = flag.String("port", "8080", "server port to bind")
		env  = flag.String("env", "develop", "exec environment (develop, production)")
	)
	flag.Parse()

	s := New()
	s.Init(*env)
	s.Middleware()
	s.Router()

	log.Printf("Starting up on http://localhost:%s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", *port), s.router))
}
