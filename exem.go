package main

import (
	"io"
	"net/http"
	"log"
	"github.com/iotdog/json2table/j2t"
)

func main() {

	handler := func(w http.ResponseWriter, req *http.Request) {
		res, _ := http.Get("http://pstgu.yss.su/iu9/mobiledev/lab4_yandex_map/?x=var09")
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()		

		_, html := j2t.JSON2HtmlTable(string(body), []string{"name", "gps", "address", "tel"}, []string{"title1"})
		log.Println(html)
		w.Write([]byte(html))
	}
	
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":300", nil))
}