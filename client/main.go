package main

import (
  "../lib/heartbeat"
  "time"
)

func main() {
  conn := &heartbeat.HeartbeatConn { IP: "127.0.0.1", Port: 1234 }
  for {
    heartbeat := heartbeat.NewHeartbeat("Client A")
    conn.Socket = heartbeat.SendHeartbeat(conn, heartbeat)
    time.Sleep(time.Second)
  }
  conn.Socket.Close()
}
