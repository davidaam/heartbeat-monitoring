package routes

import (
  "net/http"
  "github.com/gorilla/mux"
  "app/handlers"
  "app/middleware"
)

func RequireAuth(handler func(http.ResponseWriter, *http.Request)) http.Handler {
  return middleware.AuthMiddleware(http.HandlerFunc(handler))
}

func Router() *mux.Router {
  r := mux.NewRouter()

  r.Handle("/signup", http.HandlerFunc(handlers.Signup))
  r.Handle("/login", http.HandlerFunc(handlers.Login))
  r.Handle("/", RequireAuth(handlers.Index))
  r.Handle("/heartbeats", RequireAuth(handlers.Heartbeats))
  r.Handle("/heartbeats/{clientID}", RequireAuth(handlers.ClientHeartbeats))

  // Serve static files
  r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("app/src/css/"))))
  r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("app/src/js/"))))
  r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("app/src/img/"))))

  return r
}
