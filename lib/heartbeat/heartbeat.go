package heartbeat
import (
    "fmt"
    "log"
    "net"
    "strings"
    "strconv"
    "time"
)

type HeartbeatListener struct {
  IP string
  Port int64
  Timeout int32
  ReceiveCallback func(*Heartbeat)
  TimeoutCallback func(int32)
  Socket *net.UDPConn
}

type HeartbeatSender struct {
  IP string
  Port int64
  Socket net.Conn
}

type Heartbeat struct {
  ClientID string `json:"clientID"`
  Timestamp int64 `json:"timestamp"`
}

func (h *Heartbeat) toString() string {
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
    log.Panicf("[Error] %v\n", err)
  }
}

func (sender *HeartbeatSender) Send(heartbeat *Heartbeat) {
  message := heartbeat.toString()
  fmt.Printf("Sending message: %s\n", message)
  fmt.Fprintln(sender.Socket, message)
}

func (listener *HeartbeatListener) listenHeartbeats(mainChannel chan *Heartbeat) {
  buffer := make([]byte, 1024)
  for {
    n, _, err := listener.Socket.ReadFromUDP(buffer)
    checkError(err)

    heartbeat, err := parseHeartbeat(string(buffer[0:n]))
    checkError(err)

    mainChannel <- heartbeat
  }
}

func (listener *HeartbeatListener) distributeHeartbeats(mainChannel chan *Heartbeat) {
  channels := make(map[string]chan *Heartbeat)
  for {
    heartbeat := <-mainChannel
    _, channelExists := channels[heartbeat.ClientID]

    if !channelExists {
      channels[heartbeat.ClientID] = make(chan *Heartbeat)
      go listener.processHeartbeats(channels[heartbeat.ClientID])
    }

    channels[heartbeat.ClientID] <- heartbeat
  }
}

func (listener *HeartbeatListener) processHeartbeats(channel chan *Heartbeat) {
  timeouts := 0
  for {
    select {
    case heartbeat := <-channel:
      listener.ReceiveCallback(heartbeat)
      timeouts = 0
    case <-time.After(time.Second * time.Duration(listener.Timeout)):
      timeouts++
      listener.TimeoutCallback(listener.Timeout)
    }
  }
}

func (listener *HeartbeatListener) StartServer() {
  listenAddr, err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d", listener.Port))
	socket, err := net.ListenUDP("udp", listenAddr)
	checkError(err)
  defer socket.Close()

  listener.Socket = socket

  mainChannel := make(chan *Heartbeat)
  go listener.listenHeartbeats(mainChannel)
  listener.distributeHeartbeats(mainChannel)
}
