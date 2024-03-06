package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "log"
  "net"
  "os"
  "os/signal"
  "time"
)

type Peer struct {
  IP       string
  Port     int
  IsHead   bool
  NextPeer *Peer
}

type Message struct {
  Sender string
  Data   string
}

func main() {
  port := getPort()

  peer1 := &Peer{IP: "185.139.70.64", Port: 1111, IsHead: true}
  peer2 := &Peer{IP: "185.104.249.105", Port: 2222, IsHead: false}

  ip := getIPAddress()
  currentPeer := peer1

  switch ip {
  case peer1.IP:
    currentPeer = peer1
  case peer2.IP:
    currentPeer = peer2
  }

  peer1.Port = port
  peer2.Port = port

  peer1.NextPeer = peer2
  peer2.NextPeer = peer1

  go startPeer(currentPeer)
  go startCommandInterface(currentPeer)

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  <-c

  log.Println("Exiting...")
}

func getIPAddress() string {
  conn, err := net.Dial("udp", "8.8.8.8:80")
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  localAddr := conn.LocalAddr().(*net.UDPAddr)

  return localAddr.IP.String()
}

func getPort() int {
  port := 2222
  fmt.Print("Port: ")
  fmt.Scan(&port)

  address := fmt.Sprintf("localhost:%d", port)
  l, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatal(err)
  }
  defer l.Close()
  localAddr := l.Addr().(*net.TCPAddr)

  return localAddr.Port
}

func startPeer(p *Peer) {
  address := fmt.Sprintf("%s:%d", p.IP, p.Port)
  listener, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatalf("Peer %s failed to start: %v\n", address, err)
  }
  defer listener.Close()

  log.Printf("Peer started and listening on %s\n", address)

  for {
    conn, err := listener.Accept()
    if err != nil {
      log.Printf("Peer %s encountered an error: %v\n", address, err)
      continue
    }
    go handleConnection(conn, p)
  }
}

func handleConnection(conn net.Conn, p *Peer) {
  defer conn.Close()
  message := make([]byte, 4096)
  n, err := conn.Read(message)
  if err != nil {
    log.Println(err)
    return
  }
  var record Message
  err = json.Unmarshal(message[:n], &record)

  if err != nil {
    log.Println(err)
    return
  }

  if p.IP == record.Sender {
    fmt.Printf("END: %s\n", record.Data)
  } else {
    time.Sleep(1 * time.Second)

    fmt.Printf("TRANSMISSION: %s\n", record.Data)
    sendMessage(p, record)
  }
}

func startCommandInterface(p *Peer) {
  for {
    var command string
    fmt.Print("Enter command (start, stop))\n")
    fmt.Scanln(&command)
    scanner := bufio.NewScanner(os.Stdin)

    switch command {
    case "start":
      fmt.Printf("Enter message to send: ")
      scanner.Scan()
      message := scanner.Text()

      msg := Message{Sender: p.IP, Data: message}

      sendMessage(p, msg)
    case "stop":
      return
    default:
      log.Println("Invalid command")
    }
  }
}

func sendMessage(p *Peer, message Message) {
  address := fmt.Sprintf("%s:%d", p.NextPeer.IP, p.NextPeer.Port)
  conn, err := net.Dial("tcp", address)
  if err != nil {
    log.Println(err)
    return
  }
  defer conn.Close()

  jsonData, err := json.Marshal(message)
  if err != nil {
    log.Println(err)
    return
  }

  _, err = conn.Write(jsonData)
  if err != nil {
    log.Println(err)
    return
  }

  log.Printf("Sent message to neighbor %s\n", address)
}
