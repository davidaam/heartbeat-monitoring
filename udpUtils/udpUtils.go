package udpUtils
import (
    "fmt"
    "log"
    "net"
)

type UDPConn struct {
  IP string
  Port int
  Socket net.Conn
  PacketsSent int
  PacketsReceived int
}

func checkError(err error) {
  if err != nil {
    log.Fatalln("[Error] %v", err)
  }
}

func SendMessage(conn *UDPConn, message string) (net.Conn) {
  var err error
  socket := conn.Socket

  if socket == nil {
    server := fmt.Sprintf("%s:%d", conn.IP, conn.Port)
    socket, err = net.Dial("udp", server)
    checkError(err)
  }

  fmt.Printf("Sending message: %s\n", message)
  fmt.Fprintln(socket, message) // send heartbeat through socket
  conn.PacketsSent++

  return socket
}

func SendHeartbeat(conn *UDPConn) (net.Conn) {
  return SendMessage(conn, "1")
}

func ListenHeartbeats(listenConn *UDPConn) {
  listenAddr, err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d", listenConn.Port))
	listenSocket, err := net.ListenUDP("udp", listenAddr)
	checkError(err)

  buffer := make([]byte, 1024)
  for {
    n, sentAddr, err := listenSocket.ReadFromUDP(buffer)
    checkError(err)
    listenConn.PacketsReceived++
    fmt.Println("Received ",string(buffer[0:n]), " from ", sentAddr)
  }
}
