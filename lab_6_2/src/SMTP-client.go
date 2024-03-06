package main

import (
    "crypto/tls"
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "net/smtp"
    "time"

  _ "github.com/go-sql-driver/mysql"
)

type EmailConfiguration struct {
    Username string
    Password string
    Host     string
    Port     string
}

type EmailData struct {
    To      string
    Subject string
    Body    string
}

var emailConfig = EmailConfiguration{
    Username: "dts21@dactyl.su",
    Password: "12345678990DactylSUDTS",
    Host:     "mail.nic.ru",
    Port:     "465",
}

func main() {
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS lab7Gorbunov (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username TEXT NOT NULL,
        email TEXT NOT NULL,
        message TEXT NOT NULL
    );`
    db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(createTableQuery)
    if err != nil {
        log.Fatal(err)
    }

    for i := 0; i < 10; i++{
    rows, err := db.Query("SELECT username, email, message FROM lab7Gorbunov")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var username, email, message string


        err := rows.Scan(&username, &email, &message)
        if err != nil {
          log.Fatal(err)
        }

        emailData := generateEmailData(username, email, message)

        sendEmail(emailData)

        time.Sleep(time.Duration(randomDelay()) * time.Second)
    }
  }
}

func generateEmailData(username, to, message string) EmailData {
    subject := "лаб. 6_2 Горбунов Артём ИУ9-31Б"

    body := fmt.Sprintf(`
    <html>
        <head></head>
        <body>
            <p><strong>Здравствуйте, %s</strong>,</p>
            <p><em>%s</em></p>
        </body>
    </html>`, username, message)

    return EmailData{
        To:      to,
        Subject: subject,
        Body:    body,
    }
}

func sendEmail(emailData EmailData) {
    tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         emailConfig.Host,
    }

    conn, err := tls.Dial("tcp", emailConfig.Host+":"+emailConfig.Port, tlsConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)

    client, err := smtp.NewClient(conn, emailConfig.Host)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    if err = client.Auth(auth); err != nil {
        log.Fatal(err)
    }

    if err = client.Mail(emailConfig.Username); err != nil {
        log.Fatal(err)
    }

    if err = client.Rcpt(emailData.To); err != nil {
        log.Fatal(err)
    }

    dataWriter, err := client.Data()
    if err != nil {
        log.Fatal(err)
    }

    _, err = dataWriter.Write([]byte(
        "To: " + emailData.To + "\r\n" +
        "Subject: " + emailData.Subject + "\r\n" +
        "Content-Type: text/html; charset=UTF-8\r\n\r\n" +
        emailData.Body))
    if err != nil {
        log.Fatal(err)
    }

    err = dataWriter.Close()
    if err != nil {
        log.Fatal(err)
    }
}

func randomDelay() int {
    return rand.Intn(10) + 1
}
