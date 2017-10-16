package routes

import (
  "github.com/gorilla/mux"
  "../handlers"
)

func Router() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/", handlers.Heartbeats)
  r.HandleFunc("/{clientID}", handlers.ClientHeartbeats)
  return r
}
