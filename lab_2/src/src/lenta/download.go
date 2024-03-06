package main

import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
	"fmt"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Time, Title string
}

func readItem(item *html.Node) *Item {
	if a := item.FirstChild; isElem(a, "a") {
		if cs := getChildren(a); len(cs) == 2 && isElem(cs[0], "time") && isText(cs[1]) {
			return &Item{
				Ref:   getAttr(a, "href"),
				Time:  getAttr(cs[0], "title"),
				Title: cs[1].Data,
			}
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	if isDiv(node, "oby5F") {
		var items []*Item
		//fmt.Print(isDiv(node, "oby5F")," ", getAttr(node,"class"), "\n")
		for a := getChildren(node)[0]; a != nil; a = a.NextSibling {
			if isDiv(a, "wsXlA pZeTF SUiYd") {
				for b := getChildren(a)[0]; b != nil; b = b.NextSibling{
					if isDiv(b, "wsXlA pZeTF IS7Va") {
						for c := getChildren(b)[0]; c != nil; c = c.NextSibling{
							if isDiv(c, "wsXlA pZeTF CxVPH") {
								for d := getChildren(c)[0]; d != nil; d = d.NextSibling{
									if isDiv(d, "naAYr") {
										for e := getChildren(d)[0]; e != nil; e = e.NextSibling{
											if getAttr(e, "style") == "height:100%;max-height:48px" {
												for f := getChildren(e)[0]; f != nil; f = f.NextSibling{
													if getAttr(f, "class")=="vcSoT b6DKO jkBWH f5gWK" {
														for g := getChildren(f)[0]; g != nil; g = g.NextSibling{
															if isDiv(g, "mQ7Bh"){
																for h := getChildren(g)[0]; h != nil; h = h.NextSibling{
																	if isText(h){
																		fmt.Print(h.Data, " ", "https://www.afisha.ru" + getAttr(f,"href"), "\n\n")
																		/*items = append(items, &Item{
																			Ref:	getAttr(f,"href"),
																			Title: h.Data,
																		})*/
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func downloadNews() []*Item {
	log.Info("sending request to www.afisha.ru/msk/cinema")
	if response, err := http.Get("https://www.afisha.ru/msk/cinema/"); err != nil {
		log.Error("request to www.afisha.ru/msk/cinema failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from www.afisha.ru/msk/cinema", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from www.afisha.ru/msk/cinema", "error", err)
			} else {
				log.Info("HTML from www.afisha.ru/msk/cinema parsed successfully")
				//fmt.Print(html.Parse(response.Body))
				return search(doc)
			}
		}
	}
	return nil
}