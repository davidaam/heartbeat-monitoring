package handlers

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "../../../lib/dbLogger"
)

func success(w http.ResponseWriter) {
  w.WriteHeader(http.StatusOK)
}

func ClientHeartbeats(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  heartbeats := dbLogger.GetHeartbeats(vars["clientID"])
  json, _ := json.Marshal(heartbeats)
  fmt.Fprintf(w, "%v", string(json))
}

func Heartbeats(w http.ResponseWriter, r *http.Request) {
  heartbeats := dbLogger.ListHeartbeats()
  json, _ := json.Marshal(heartbeats)
  fmt.Fprintf(w, "%v", string(json))
}
