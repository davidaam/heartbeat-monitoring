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
  db := dbLogger.InitializeDB(sqliteDB)

  receiveCallback := func (heartbeat *heartbeat.Heartbeat) {
    fmt.Println("Heartbeat received")
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
