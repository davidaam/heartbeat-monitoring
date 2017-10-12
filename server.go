package main

import "./udpUtils"

func main() {
  udpUtils.ListenHeartbeats(&udpUtils.UDPConn { Port: 1234 })
}
