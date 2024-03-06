package main

import (
	"encoding/xml"
	"fmt" // пакет для форматированного ввода вывода
	"io/ioutil"
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	// пакет для работы с UTF-8 строками
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //анализ аргументов,

	fmt.Fprintf(w, "<p>"+r.FormValue("")+"</p>")
	fmt.Fprintf(w, "<a href='/' >главная</a><br><a href='/rss' >rss</a><br><a href='/about' >about</a>") // отправляем данные на клиентскую сторону
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://news.rambler.ru/rss/Namibia/")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	rss := RSS{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println("Заголовок канала:", rss.Channel.Title)
	fmt.Println("Описание канала:", rss.Channel.Description)
	fmt.Println("Количество элементов:", len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		fmt.Println("------------------")
		fmt.Println("Заголовок:", item.Title)
		fmt.Println("Описание:", item.Description)
		fmt.Println("Ссылка:", item.Link)
	}
	r.ParseForm()
	fmt.Fprintf(w, "<p>Заголовок канала:"+rss.Channel.Title+"</p>")
	fmt.Fprintf(w, "<p>Описание канала:"+rss.Channel.Description+"</p>")
	for _, item := range rss.Channel.Items {
		fmt.Fprintf(w, "<p>------------------</p>")
		fmt.Fprintf(w, "<p>Заголовок:"+item.Title+"</p>")
		fmt.Fprintf(w, "<p>Описание:"+item.Description+"</p>")
		fmt.Fprintf(w, "<p>Ссылка:"+item.Link+"</p>")
	}
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //анализ аргументов,
	fmt.Fprintf(w, "<p>"+r.FormValue("")+"</p>")
	fmt.Fprintf(w, "<a href='/' >главная</a><br>about") // отправляем данные на клиентскую сторону
}

func main() {
	http.HandleFunc("/", HomeRouterHandler) // установим роутер
	http.HandleFunc("/rss", rssHandler)
	http.HandleFunc("/about", aboutHandler)
	err := http.ListenAndServe(":9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
