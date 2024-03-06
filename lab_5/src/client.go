package main

import (
  "encoding/json"
  "fmt"
  "log"

  "github.com/gorilla/websocket"
  "github.com/skorobogatov/input"
)

var addr = "" // адрес сервера WebSocket

func sendMessageToServer(numbers []float64) ([]byte, error) {
  // Устанавливаем соединение с сервером WebSocket
  c, _, err := websocket.DefaultDialer.Dial(addr, nil)
  if err != nil {
    return nil, err
  }
  defer c.Close()

  // Подготовка данных для отправки на сервер
  requestBody := map[string]interface{}{"numbers": numbers}
  jsonBody, err := json.Marshal(requestBody)
  if err != nil {
    return nil, err
  }

  // Отправляем сообщение на сервер
  err = c.WriteMessage(websocket.TextMessage, jsonBody)
  if err != nil {
    return nil, err
  }

  // Получаем ответ от сервера
  _, response, err := c.ReadMessage()
  if err != nil {
    return nil, err
  }

  return response, nil
}

func main() {
    fmt.Print("IP: ")
    ip := ""
    fmt.Scan(&ip)
    fmt.Print("Port: ")
    port := ""
    fmt.Scan(&port)
    addr = "ws://" + ip + ":" + port
  for {
    fmt.Println("Введите команду ('quit' для выхода, 'calc' для отправки данных):")
    command := input.Gets()

    switch command {
    case "quit":
      return
    case "calc":
      var numbers []float64
      var n int
      fmt.Print("Введите количество чисел: ")
      fmt.Scanln(&n)
      fmt.Println("Введите числа:")
      for i := 0; i < n; i++ {
        var num float64
        fmt.Scanln(&num)
        numbers = append(numbers, num)
      }

      response, err := sendMessageToServer(numbers)
      if err != nil {
        log.Fatal("Ошибка при отправке сообщения на сервер:", err)
      }

      fmt.Println("Ответ от сервера:", string(response))
    default:
      fmt.Println("Ошибка: неизвестная команда")
    }
  }
}