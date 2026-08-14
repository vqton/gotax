package web

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// legacyPage is the data for rendering an unconverted Alpine page inside the
// shared shell. BodyAttrs/Content are lifted from the static file: the page's
// own <body> attributes (x-data, x-init) and its inner content minus the
// duplicate shell containers (aside#sidebar, header#topbar) and the
// <div class="lg:ml-60"><main> wrappers, which the server shell replaces.
type legacyPage struct {
	Title   string
	Content template.HTML
}

// loadLegacy reads a legacy static page from web/app and extracts the parts
// that belong in the shell's content block. Kept as plain HTML so Alpine
// pages keep working unchanged; the extraction is structure-driven and
// relies on the uniform legacy boilerplate (verified across all pages).
func loadLegacy(path string) (legacyPage, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return legacyPage{}, err
	}
	doc, err := html.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return legacyPage{}, fmt.Errorf("parse %s: %w", path, err)
	}

	body := findTag(doc, "body")
	if body == nil {
		return legacyPage{}, fmt.Errorf("%s: no body", path)
	}
	title := findTitle(doc)

	// Remove duplicate shell containers; the server-rendered shell replaces
	// them. Direct children only — nested containers (topbar inside the
	// lg:ml-60 wrapper) are dropped by the unwrap below.
	for _, id := range []string{"sidebar", "topbar"} {
		if n := findByID(body, id); n != nil && n.Parent == body {
			body.RemoveChild(n)
		}
	}
	// Unwrap <div class="lg:ml-60"> and <main class="p-6"> so content sits
	// directly in the shell's content area.
	for _, sel := range []string{"lg:ml-60"} {
		if n := findClass(body, sel); n != nil && n.Parent != nil {
			unwrap(n)
		}
	}
	if n := firstMain(body); n != nil && n.Parent != nil {
		unwrap(n)
	}

	// Wrap remaining content in a div carrying the original body attributes
	// (x-data, x-init, class) so the Alpine component still boots.
	wrapper := &html.Node{Type: html.ElementNode, Data: "div", Attr: body.Attr}
	for c := body.FirstChild; c != nil; {
		next := c.NextSibling
		body.RemoveChild(c)
		wrapper.AppendChild(c)
		c = next
	}

	var sb strings.Builder
	if err := html.Render(&sb, wrapper); err != nil {
		return legacyPage{}, fmt.Errorf("render %s: %w", path, err)
	}
	title = strings.TrimSuffix(strings.TrimSpace(title), " · GoTax")
	return legacyPage{Title: title, Content: template.HTML(sb.String())}, nil
}

func findTag(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
		if r := findTag(c, tag); r != nil {
			return r
		}
	}
	return nil
}

func findByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if r := findByID(c, id); r != nil {
			return r
		}
	}
	return nil
}

func findClass(n *html.Node, class string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, class) {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if r := findClass(c, class); r != nil {
			return r
		}
	}
	return nil
}

func firstMain(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "main" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if r := firstMain(c); r != nil {
			return r
		}
	}
	return nil
}

// unwrap replaces n with its own children.
func unwrap(n *html.Node) {
	parent := n.Parent
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		n.RemoveChild(c)
		parent.InsertBefore(c, n)
		c = next
	}
	parent.RemoveChild(n)
}

func findTitle(doc *html.Node) string {
	head := findTag(doc, "head")
	if head == nil {
		return ""
	}
	if t := findTag(head, "title"); t != nil && t.FirstChild != nil {
		return strings.TrimSpace(t.FirstChild.Data)
	}
	return ""
}
