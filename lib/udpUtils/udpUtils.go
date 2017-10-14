package udpUtils
import (
    "fmt"
    "log"
    "net"
    "strings"
    "strconv"
    "time"
)

type HeartbeatConn struct {
  IP string
  Port int64
  Socket net.Conn
  PacketsSent int64
  PacketsReceived int64
  LastHeartbeatTime int64
  Alive bool
}

type Heartbeat struct {
  ClientID string
  Timestamp int64
}

func (h *Heartbeat) to_string() string {
  return fmt.Sprintf("%s|%d", h.ClientID, h.Timestamp)
}

func NewHeartbeat (clientID string) *Heartbeat {
  return &Heartbeat { ClientID: clientID, Timestamp: time.Now().Unix() }
}

func parseHeartbeat (message string) (*Heartbeat, error) {
  arr := strings.Split(message, "|")
  clientId := arr[0]
  timestamp, err := strconv.ParseInt(strings.TrimSpace(arr[1]), 10, 64)
  return &Heartbeat{ ClientID: clientId, Timestamp: timestamp }, err
}

func checkError(err error) {
  if err != nil {
    log.Fatalln("[Error] %v", err)
  }
}

func SendMessage(conn *HeartbeatConn, message string) (net.Conn) {
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

func SendHeartbeat(conn *HeartbeatConn, heartbeat *Heartbeat) (net.Conn) {
  return SendMessage(conn, heartbeat.to_string())
}

func listenHeartbeats(listenSocket *net.UDPConn, listenConn *HeartbeatConn, receiveCallback func(*Heartbeat)) {
  buffer := make([]byte, 1024)
  for {
    n, _, err := listenSocket.ReadFromUDP(buffer)
    checkError(err)

    heartbeat, err := parseHeartbeat(string(buffer[0:n]))
    checkError(err)
    listenConn.LastHeartbeatTime = heartbeat.Timestamp

    listenConn.PacketsReceived++
    listenConn.Alive = true
    receiveCallback(heartbeat)
  }
}

func StartServer(listenConn *HeartbeatConn, receiveCallback func(*Heartbeat), noHeartbeatsCallback func(*HeartbeatConn)) {
  listenAddr, err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d", listenConn.Port))
	listenSocket, err := net.ListenUDP("udp", listenAddr)
	checkError(err)

  go listenHeartbeats(listenSocket, listenConn, receiveCallback)

  for {
    time.Sleep(time.Second)
    // If no heartbeat has been sent ever, the noHeartbeatsCallback won't be called, since it's
    // waiting for the client to initially launch.
    if listenConn.Alive && listenConn.LastHeartbeatTime > 60 && time.Now().Unix() - listenConn.LastHeartbeatTime > 60 {
      noHeartbeatsCallback(listenConn)
      listenConn.Alive = false
    }
  }
}
