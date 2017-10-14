package main

import (
  "./lib/udpUtils"
  "time"
)

func main() {
  conn := &udpUtils.HeartbeatConn { IP: "127.0.0.1", Port: 1234 }
  for {
    heartbeat := udpUtils.NewHeartbeat("Client A")
    conn.Socket = udpUtils.SendHeartbeat(conn, heartbeat)
    time.Sleep(time.Second)
  }
  conn.Socket.Close()
}
