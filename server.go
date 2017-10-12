package main

import (
  "./udpUtils"
  "fmt"
  "time"
)

func main() {
  receiveCallback := func (message string) {
    // TODO: Everytime a heartbeat is received, log in database
    fmt.Println("Heartbeat received")
  }
  noHeartBeatsCallback := func (conn *udpUtils.HeartBeatConn) {
    fmt.Printf("No heartbeat received in the last %d seconds...\n", time.Now().Unix() - conn.LastHeartBeatTime)
  }
  udpUtils.StartServer(&udpUtils.HeartBeatConn { Port: 1234 }, receiveCallback, noHeartBeatsCallback)
}
