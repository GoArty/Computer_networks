package main

import (
	"github.com/mgutz/logxi/v1"
	"html/template"
	"net/http"
)
//var 11
const INDEX_HTML = `
    <!doctype html>
    <html lang="ru">
        <head>
            <meta charset="utf-8">
            <title>Cписка фильмов с www.afisha.ru/msk/cinema</title>
            <style>
            a {
                color: black;
                text-decoration: none;
                background-color: yellow;
                border: solid;
                border-radius: 5px;
                text-alight: center;
                hight: 50px;
                weight: 100px;
                font-family: 'Bradley Hand', cursive;
                padding: 5px 10px 5px 10px;
            }    
            ul{
                list-style-type: none;
            }
            body{
                background-color: grey;
            }

            </style>
        </head>
        <body>
            <ul>
            {{range .}}
                <li><a href="https://www.afisha.ru{{.Ref}}">{{.Title}}</a></li>
                <br>
            {{end}}
            </ul>
        </body>
    </html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, downloadNews()); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:6060", nil))
}
