package main

import (
  "../lib/heartbeat"
  "../lib/dbLogger"
  "fmt"
  "path/filepath"
)

func main() {
  sqliteDB, _ := filepath.Abs("heartbeats.sqlite3")

  fmt.Printf("*** Heartbeats will be logged in %s ***\n\n", sqliteDB)
  db := dbLogger.InitializeDB(fmt.Sprintf("file:%s?cache=shared&mode=rwc", sqliteDB))

  receiveCallback := func (heartbeat *heartbeat.Heartbeat) {
    fmt.Printf("Heartbeat received from %s\n", heartbeat.ClientID)
    go dbLogger.LogHeartbeat(db, heartbeat)
  }
  timeoutCallback := func (timeout int32) {
    fmt.Printf("No heartbeat received in %d seconds...\n", timeout)
  }

  listener := heartbeat.HeartbeatListener {
    Port: 1234,
    ReceiveCallback: receiveCallback,
    TimeoutCallback: timeoutCallback,
    Timeout: 10,
  }

  listener.StartServer()
}
