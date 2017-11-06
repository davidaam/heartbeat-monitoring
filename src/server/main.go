package main

import (
  "lib/heartbeat"
  "lib/dbLogger"
  "fmt"
  "path/filepath"
)

func main() {
  sqliteDB, _ := filepath.Abs("db.sqlite3")

  fmt.Printf("*** Heartbeats will be logged in %s ***\n\n", sqliteDB)
  db := dbLogger.InitializeDB(fmt.Sprintf("file:%s?cache=shared&mode=rwc", sqliteDB))

  heartbeats := make(chan *heartbeat.Heartbeat)
  go dbLogger.LogHeartbeats(db, heartbeats)

  receiveCallback := func (heartbeat *heartbeat.Heartbeat) {
    fmt.Printf("Heartbeat received from %s\n", heartbeat.ClientID)
    heartbeats <- heartbeat
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

  catchPanic := func() {
    f := func() {
        if r := recover(); r != nil {
        fmt.Println(r)
        fmt.Println("Recovering...")
        listener.StartServer()
      }
    }
    defer f()
  }

  listener.StartServer()
  defer catchPanic()
}
