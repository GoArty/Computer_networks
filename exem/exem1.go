package main

import (
	//"fmt"
	"io"
	"log"
	"net/http"
	//"strings"
	//"encoding/json"
	"github.com/iotdog/json2table/j2t"

)

func main() {

	handler := func(w http.ResponseWriter, req *http.Request) {
		res, err := http.Get("http://pstgu.yss.su/iu9/mobiledev/lab4_yandex_map/?x=var09")
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()		

		ok, html := j2t.JSON2HtmlTable(string(body), []string{"name", "gps", "address", "tel"}, []string{"title1"})
		if ok {
			log.Println(html)
			w.Write([]byte(html))
		} else {
			log.Println("failed to convert json to html table")
		}
	}
	
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}