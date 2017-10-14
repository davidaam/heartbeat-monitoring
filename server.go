package main

import (
  "./lib/udpUtils"
  "./lib/dbLogger"
  "fmt"
  "time"
)

func main() {
  db := dbLogger.InitializeDB("heartbeats.sqlite3")

  receiveCallback := func (heartbeat *udpUtils.Heartbeat) {
    fmt.Println("Heartbeat received")
    go dbLogger.LogHeartbeat(db, heartbeat)
  }
  noHeartbeatsCallback := func (conn *udpUtils.HeartbeatConn) {
    fmt.Printf("No heartbeat received in the last %d seconds...\n", time.Now().Unix() - conn.LastHeartbeatTime)
  }
  udpUtils.StartServer(&udpUtils.HeartbeatConn { Port: 1234 }, receiveCallback, noHeartbeatsCallback)
}
