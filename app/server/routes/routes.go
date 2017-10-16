package routes

import (
  "github.com/gorilla/mux"
  "../handlers"
)

func Router() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/heartbeats", handlers.Heartbeats)
  r.HandleFunc("/heartbeats/{clientID}", handlers.ClientHeartbeats)
  return r
}
