package udpUtils
import (
    "fmt"
    "log"
    "net"
    //"bufio"
)

type HeartBeatConn struct {
  IP string
  Port int
  Socket net.Conn
}

func checkError(err error) {
  if err != nil {
    log.Fatalln("[Error] %v", err)
  }
}

func SendMessage(conn *HeartBeatConn, message string) (net.Conn) {
  var err error
  socket := conn.Socket

  if socket == nil {
    server := fmt.Sprintf("%s:%d", conn.IP, conn.Port)
    socket, err = net.Dial("udp", server)
    checkError(err)
  }

  //buffer :=  make([]byte, 2048)

  fmt.Printf("Sending message: %s\n", message)
  fmt.Fprintln(socket, message) // send heartbeat through socket

  /* Not receiving ACK yet
  _, err = bufio.NewReader(socket).Read(buffer)
  checkError(err)
  */
  return socket
}

func SendHeartbeat(conn *HeartBeatConn) (net.Conn) {
  return SendMessage(conn, "1")
}

func ListenHeartbeats(port int) {
  listenAddr, err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d", port))
	socket, err := net.ListenUDP("udp", listenAddr)
	checkError(err)

  buffer := make([]byte, 1024)
  for {
      n, sentAddr, err := socket.ReadFromUDP(buffer)
      fmt.Println("Received ",string(buffer[0:n]), " from ", sentAddr)
      checkError(err)
  }
}
