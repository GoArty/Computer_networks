package main

import (
  	"bytes"
  	"database/sql"
  	"encoding/xml"
  	"fmt"
  	"golang.org/x/net/html/charset"
  	"io/ioutil"
  	"log"
  	"net/http"

  _ "github.com/go-sql-driver/mysql"
)

type Item struct {
	Guid		string `xml:"guid"`
  	Title       string `xml:"title"`
  	Link        string `xml:"link"`
  	Description string `xml:"description"`
  	PubDate     string `xml:"pub_date"`
	Enclosure   struct {
		URL  	string `xml:"url,attr"`
		Type 	string `xml:"type,attr"`
	}				   `xml:"enclosure"`
	Category	string `xml:"category"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
  	Items   []Item   `xml:"channel>item"`
}

func main() {
	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
  	if err != nil {
    	log.Fatal(err)
  	}
  	defer db.Close()

  	err = db.Ping()
  	if err != nil {
  		log.Fatal(err)
  	}

  	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS lab6Gorbunov (
      	id INT AUTO_INCREMENT PRIMARY KEY,
      	title TEXT,
      	link TEXT,
      	description TEXT,
      	pub_date DATE,
		enclosure_url TEXT,
		enclosure_type TEXT,
		category TEXT
    	)
  	`)
  	if err != nil {
    	log.Fatal(err)
  	}

  	_, err = db.Exec(`
    	ALTER TABLE lab6Gorbunov CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
  	`)
  	if err != nil {
    	log.Fatal(err)
  	}

  	resp, err := http.Get("https://vesti-k.ru/rss/")
  	if err != nil {
    	log.Fatal(err)
  	}
  	defer resp.Body.Close()

  	body, err := ioutil.ReadAll(resp.Body)
  	if err != nil {
    	log.Fatal(err)
  	}

  	decoder := xml.NewDecoder(bytes.NewReader(body))
  	decoder.CharsetReader = charset.NewReaderLabel

  	var rss RSS
  	err = decoder.Decode(&rss)
  	if err != nil {
    	log.Fatal(err)
  	}

  	for _, item := range rss.Items {
    	var count int
    	err := db.QueryRow("SELECT COUNT(*) FROM lab6Gorbunov WHERE title = ?", item.Title).Scan(&count)
    	if err != nil {
      		log.Fatal(err)
    	}

    
    	if count == 0 {
      		_, err := db.Exec("INSERT INTO lab6Gorbunov (title, link, description, pub_date, enclosure_url, enclosure_type, category) VALUES (?, ?, ?, ?, ?, ?, ?)",
        		item.Title, item.Link, item.Description, item.PubDate, item.Enclosure.URL, item.Enclosure.Type, item.Category)
      		if err != nil {
        		log.Fatal(err)
      		}
      		fmt.Printf("Добавлена новая новость: %s\n", item.Title)
    	}
  }
}
