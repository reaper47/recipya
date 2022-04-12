package scraper

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func getElement(n *html.Node, attr, id string) *html.Node {
	return traverse(n, attr, id)
}

func traverse(n *html.Node, attr, id string) *html.Node {
	if n.Type == html.ElementNode {
		s := getAttr(n, attr)
		if s == id {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, attr, id)
		if result != nil {
			return result
		}
	}
	return nil
}

func traverseAll(doc *html.Node, matcher func(node *html.Node) bool) (nodes []*html.Node) {
	var keep bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return nodes
}

func getAttr(node *html.Node, key string) string {
	if node != nil {
		for _, attr := range node.Attr {
			if attr.Key == key {
				return attr.Val
			}
		}
	}
	return ""
}

func getElementData(root *html.Node, attr, id string) chan string {
	ch := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			ch <- s
		}()

		n := getElement(root, attr, id)
		s = strings.TrimSpace(n.FirstChild.Data)
	}()
	return ch
}

func getItemPropAttr(root *html.Node, prop, attr string) chan string {
	ch := make(chan string)
	go func() {
		node := getElement(root, "itemprop", prop)
		var v string
		if node != nil {
			v = getAttr(node, attr)
		}
		ch <- v
	}()
	return ch
}

func getItemPropData(root *html.Node, item string) chan string {
	ch := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			ch <- v
		}()
		node := getElement(root, "itemprop", item)
		v = node.FirstChild.Data
	}()
	return ch
}

func findYield(parts []string) int16 {
	for _, part := range parts {
		i, err := strconv.ParseInt(part, 10, 16)
		if err == nil {
			return int16(i)
		}
	}
	return 0
}
