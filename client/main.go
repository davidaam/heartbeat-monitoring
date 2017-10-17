package main

import (
  "../lib/heartbeat"
  "time"
  "fmt"
  "log"
  "net"
  "os"
)

func main() {
  clientID := os.Args[1]

  sender := &heartbeat.HeartbeatSender { IP: "127.0.0.1", Port: 1234 }
  server := fmt.Sprintf("%s:%d", sender.IP, sender.Port)
  socket, err := net.Dial("udp", server)
  if err != nil {
    log.Fatalln(err)
  }
  sender.Socket = socket

  defer sender.Socket.Close()

  for {
    sender.Send(heartbeat.NewHeartbeat(clientID))
    time.Sleep(time.Second)
  }
}
