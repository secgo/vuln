package linkhref

import (
	"io"

	"golang.org/x/net/html"
)

type Lurl struct {
	Href string
}

func Parse(r io.Reader, tag, attribut string) ([]Lurl, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := LinkNodes(doc, tag)
	var Link []Lurl
	for _, node := range nodes {
		Link = append(Link, buildLink(node, attribut))
	}
	return Link, nil
}

func LinkNodes(n *html.Node, tag string) []*html.Node {
	var r []*html.Node
	if n.Type == html.ElementNode && n.Data == tag {
		return []*html.Node{n}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r = append(r, LinkNodes(c, tag)...)
	}
	return r
}

func buildLink(n *html.Node, attribute string) Lurl {
	var ret Lurl
	for _, attr := range n.Attr {
		if attr.Key == attribute {
			ret.Href = attr.Val
			break
		}
	}
	return ret
}
