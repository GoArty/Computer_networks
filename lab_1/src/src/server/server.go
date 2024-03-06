package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"proto"
	"strings"

	log "github.com/mgutz/logxi/v1"
)

// Функция для вычисления значения полинома в точке x0.
func evaluatePolynomial(str string, substr string) []int {
	indices := []int{}
	startIndex := 0

	// Ищем подстроку в строке, пока не достигнем конца строки
	for {
		index := strings.Index(str[startIndex:], substr)

		// Если подстрока не найдена, прекращаем поиск
		if index == -1 {
			break
		}

		// Добавляем индекс найденной подстроки
		indices = append(indices, startIndex+index)

		// Увеличиваем стартовый индекс для следующего поиска
		startIndex += index + len(substr)
	}

	return indices
}

// Client - состояние клиента.
type Client struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений
	stro   proto.Stroca
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		stro:   proto.Stroca{Stroca: ""},
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)

	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			if client.handleRequest(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "stroca":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			if err := json.Unmarshal(*req.Data, &client.stro); err != nil {
				errorMsg = "malformed data field"
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "substring":
		if client.stro.Stroca == "" {
			client.respond("failed", "polynomial coefficients are not set")
		} else if req.Data == nil {
			client.respond("failed", "data field is absent")
		} else {
			var str string
			if err := json.Unmarshal(*req.Data, &str); err != nil {
				client.respond("failed", "malformed data field")
			} else {
				result := evaluatePolynomial(client.stro.Stroca, str)
				client.respond("result", &result)
			}
		}

	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{Status: status, Data: &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())
		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn).serve()
				}
			}
		}
	}
}
