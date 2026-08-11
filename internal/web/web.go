// Package web serves the server-rendered htmx frontend (replaces the
// Alpine.js client-rendered pages). Each page has its own template set:
// base/sidebar/topbar partials + the page's own template, parsed together so
// every page may define "content" without clashing.
package web

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	baseDir   = "web/templates"
	pagesDir  = "web/app"
	payDir    = "web/payroll"
	baseFiles = "web/templates/base.html,web/templates/_sidebar.html,web/templates/_topbar.html"
)

func init() { _ = baseFiles }

// Server renders pages from pre-parsed template sets.
type Server struct {
	sets map[string]*template.Template
}

// NewServer parses every page template into its own set. pages lists the
// page names (route keys, e.g. "dashboard") whose templates exist in web/app.
func NewServer(pages []string) (*Server, error) {
	partials := []string{
		filepath.Join(baseDir, "base.html"),
		filepath.Join(baseDir, "_sidebar.html"),
		filepath.Join(baseDir, "_topbar.html"),
	}
	s := &Server{sets: map[string]*template.Template{}}
	for _, p := range pages {
		files := append(append([]string{}, partials...), pageFile(p))
		if _, err := os.Stat(files[len(files)-1]); err != nil {
			return nil, fmt.Errorf("page %s: %w", p, err)
		}
		tmpl, err := template.New("base").Funcs(funcs).ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("parse page %s: %w", p, err)
		}
		s.sets[p] = tmpl
	}
	return s, nil
}

func pageFile(page string) string {
	if strings.HasPrefix(page, "payroll-") {
		return filepath.Join(payDir, strings.TrimPrefix(page, "payroll-")+".html")
	}
	return filepath.Join(pagesDir, page+".html")
}

// Render executes the "base" template of a page set. data is the page data;
// the layout receives {Shell, Data}.
func (s *Server) Render(c *gin.Context, page string, shell Shell, data any) {
	set, ok := s.sets[page]
	if !ok {
		http.Error(c.Writer, "page template not found: "+page, http.StatusNotFound)
		return
	}
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	if err := set.ExecuteTemplate(c.Writer, "base", gin.H{"Shell": shell, "Data": data}); err != nil {
		log.Printf("render page %s: %v", page, err)
	}
}

// RenderFragment executes a named fragment (e.g. "users-table") of a page
// set — used for htmx swaps after mutations. Data is wrapped in a "Data" key
// so fragments read `.Data.X` identically to full page renders.
func (s *Server) RenderFragment(c *gin.Context, page, frag string, data any) {
	set, ok := s.sets[page]
	if !ok {
		http.Error(c.Writer, "page template not found: "+page, http.StatusNotFound)
		return
	}
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	if err := set.ExecuteTemplate(c.Writer, frag, gin.H{"Data": data}); err != nil {
		log.Printf("render fragment %s/%s: %v", page, frag, err)
	}
}

// Toast sets the HX-Trigger header so the client shows a toast after an htmx
// mutation. Must be set before writing the response body.
func Toast(c *gin.Context, typ, text string) {
	c.Header("HX-Trigger", fmt.Sprintf(`{"toast":{"type":%q,"text":%q}}`, typ, text))
}

// ErrPageLoad is wrapped by page loaders when data gathering fails.
var ErrPageLoad = errors.New("page load failed")

/* ─── Template functions ─── */

var funcs = template.FuncMap{
	"fmtVND": func(v any) string {
		var f float64
		switch n := v.(type) {
		case int:
			f = float64(n)
		case int64:
			f = float64(n)
		case float64:
			f = n
		case float32:
			f = float64(n)
		default:
			return "0"
		}
		return formatVND(f)
	},
	"fmtDate":     fmtDate,
	"fmtDateTime": fmtDateTime,
	"today": func() string {
		return time.Now().Format("2006-01-02")
	},
	"roleLabel":   roleLabel,
	"statusBadge": statusBadge,
	"avatar":      avatarInit,
	"list": func(items ...any) []any {
		return items
	},
	"tup": func(a, b, c string) []string {
		return []string{a, b, c}
	},
	"dict": func(kv ...any) map[string]any {
		m := map[string]any{}
		for i := 0; i+1 < len(kv); i += 2 {
			m[fmt.Sprint(kv[i])] = kv[i+1]
		}
		return m
	},
	"string": func(v any) string { return fmt.Sprint(v) },
}
