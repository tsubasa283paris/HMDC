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
	h := api.NewHandler()
	s.router.Route("/api", func(api chi.Router) {
		api.Use(Auth("db connection")) // middleware

		// users API
		api.Route("/users", func(users chi.Router) {
			users.Get("/", h.GetUsers)
		})
	})

	// auth API
	s.router.Route("/api/auth", func(auth chi.Router) {
		auth.Use(Check)
		auth.Post("/login", h.Login)
	})
}

// Authentication
func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "admin" {
				api.RespondError(w, "invalid token", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	return
}

func Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// body, _ := io.ReadAll(r.Body)
		// r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
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
