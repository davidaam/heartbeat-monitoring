package main

import (
  "../lib/heartbeat"
  "../lib/dbLogger"
  "fmt"
  "time"
)

func main() {
  db := dbLogger.InitializeDB("heartbeats.sqlite3")

  receiveCallback := func (heartbeat *heartbeat.Heartbeat) {
    fmt.Println("Heartbeat received")
    go dbLogger.LogHeartbeat(db, heartbeat)
  }
  noHeartbeatsCallback := func (conn *heartbeat.HeartbeatConn) {
    fmt.Printf("No heartbeat received in the last %d seconds...\n", time.Now().Unix() - conn.LastHeartbeatTime)
  }
  heartbeat.StartServer(&heartbeat.HeartbeatConn { Port: 1234 }, receiveCallback, noHeartbeatsCallback)
}
