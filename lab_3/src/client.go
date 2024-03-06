package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "time"
   )

type AddressBook struct {
 Records []Record `json:"records"`
}

type Record struct {
 Name  string `json:"name"`
 Email string `json:"email"`
}
type SupType struct{
    supType ForJSON `json:"suptype"`
}
type ForJSON struct{
    forJSON Record `json:"forjson"`
    oper string `json:"oper"`
    index int `json:"index"`
    len int `json:"len"`
 }

type Peer struct {
 IP   string
 Port int
}

func main() {
 ip := getIPAddress()
 port := getPort()
 neighbors := getNeighbors()
 addressBook := AddressBook{}

 go startServer(ip, port, neighbors, &addressBook)

 go startCommandServer(&addressBook, port)

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

func getNeighbors() []Peer {
 neighbors := []Peer{{"185.139.70.64", 1},{"185.104.249.105", 2},{"185.255.133.113", 3}}
 return neighbors
}

func startServer(ip string, port int, neighbors []Peer, addressBook *AddressBook) {
 address := fmt.Sprintf("%s:%d", ip, port)
 listener, err := net.Listen("tcp", address)
 if err != nil {
  log.Fatal(err)
 }
 defer listener.Close()
 log.Printf("Server started and listening on %s\n", address)
 for {
  conn, err := listener.Accept()
  if err != nil {
   log.Fatal(err)
  }
  go handleConnection(conn, addressBook)
 }
}

func handleConnection(conn net.Conn, addressBook *AddressBook) {
 defer conn.Close()
 log.Println("New connection established")
 message := make([]byte, 4096)
 n, err := conn.Read(message)
 if err != nil {
  log.Println(err)
  return
 }
 var record AddressBook
 err = json.Unmarshal(message[:n], &record)
 
 if err != nil {
  log.Println(err)
  return
 }
 addressBook.Records = record.Records 
}

func sendAddressBookToNeighbors(addressBook *AddressBook, port int) {
    jsonData, err := json.Marshal(addressBook)
    if err != nil {
     log.Println(err)
     return
    }
    for _, neighbor := range getNeighbors() {
      if(neighbor != ([]Peer{{getIPAddress(), port}})[0]){
         go sendAddressBook(neighbor, jsonData)
      }
    }
}

func sendAddressBook(neighbor Peer, jsonData []byte) {
 address := fmt.Sprintf("%s:%d", neighbor.IP, neighbor.Port)
 conn, err := net.DialTimeout("tcp", address, time.Second*5)
 if err != nil {
  log.Printf("Failed to connect to neighbor %s: %v\n", address, err)
  return
 }
 defer conn.Close()

 _, err = conn.Write(jsonData)
 if err != nil {
  log.Println(err)
  return
 }

 log.Printf("Sent address book to neighbor %s\n", address)
}

func startCommandServer(addressBook *AddressBook, port int) {
 for {
  var command string
  fmt.Print("Enter command (add, remove, list): ")
  fmt.Scanln(&command)

  switch command {
  case "add":
   addRecord(addressBook, port)
  case "remove":
   removeRecord(addressBook, port)
  case "list":
   listRecords(addressBook)
  default:
   log.Println("Invalid command")
}
}
}

func addRecord(addressBook *AddressBook, port int) {
var name, email string
fmt.Print("Enter name: ")
fmt.Scanln(&name)
fmt.Print("Enter email: ")
fmt.Scanln(&email)

record := Record{Name: name, Email: email}
addressBook.Records = append(addressBook.Records, record)

sendAddressBookToNeighbors(addressBook, port)
}

func removeRecord(addressBook *AddressBook, port int) {
var index int
fmt.Print("Enter the index of the record to remove: ")
fmt.Scanln(&index)

if index >= 0 && index < len(addressBook.Records) {
addressBook.Records = append(addressBook.Records[:index], addressBook.Records[index+1:]...)

sendAddressBookToNeighbors(addressBook, port)
} else {
log.Println("Invalid index")
}
}

func listRecords(addressBook *AddressBook) {
for i, record := range addressBook.Records {
fmt.Printf("%d. %s - %s\n", i, record.Name, record.Email)
}
}
