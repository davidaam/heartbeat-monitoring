package main

import (
  "./lib/udpUtils"
  "time"
)

func main() {
  conn := udpUtils.HeartbeatConn { IP: "127.0.0.1", Port: 1234 }
  connPointer := &conn
  for {
    connPointer.Socket = udpUtils.SendHeartbeat(connPointer)
    time.Sleep(time.Second)
  }
  conn.Socket.Close()
}
