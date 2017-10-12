package udpUtils
import (
    "fmt"
    "log"
    "net"
    "strings"
    "strconv"
    "time"
)

type HeartBeatConn struct {
  IP string
  Port int64
  Socket net.Conn
  PacketsSent int64
  PacketsReceived int64
  LastHeartBeatTime int64
  Alive bool
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

  fmt.Printf("Sending message: %s\n", message)
  fmt.Fprintln(socket, message) // send heartbeat through socket
  conn.PacketsSent++

  return socket
}

func SendHeartbeat(conn *HeartBeatConn) (net.Conn) {
  return SendMessage(conn, strconv.Itoa(int(time.Now().Unix())))
}

func listenHeartbeats(listenSocket *net.UDPConn, listenConn *HeartBeatConn, receiveCallback func(string)) {
  buffer := make([]byte, 1024)
  for {
    n, _, err := listenSocket.ReadFromUDP(buffer)
    checkError(err)

    receivedMessage := string(buffer[0:n])
    listenConn.LastHeartBeatTime, err = strconv.ParseInt(strings.TrimSpace(receivedMessage), 10, 64)
    checkError(err)

    listenConn.PacketsReceived++
    listenConn.Alive = true
    receiveCallback(receivedMessage)
  }
}

func StartServer(listenConn *HeartBeatConn, receiveCallback func(string), noHeartBeatsCallback func(*HeartBeatConn)) {
  listenAddr, err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d", listenConn.Port))
	listenSocket, err := net.ListenUDP("udp", listenAddr)
	checkError(err)

  go listenHeartbeats(listenSocket, listenConn, receiveCallback)

  for {
    time.Sleep(time.Second)
    // If no heartbeat has been sent ever, the noHeartBeatsCallback won't be called, since it's
    // waiting for the client to initially launch.
    if listenConn.Alive && listenConn.LastHeartBeatTime > 60 && time.Now().Unix() - listenConn.LastHeartBeatTime > 60 {
      noHeartBeatsCallback(listenConn)
      listenConn.Alive = false
    }
  }
}
