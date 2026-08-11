package web

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"gotax/internal/auth"
)

// PageAuthMiddleware authenticates page requests from the Authorization
// header (htmx requests, set by app.js) or the gotax_token cookie (full page
// loads, set by GoTax.Auth.saveLogin). Unauthenticated → redirect /login.
func PageAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
			token = strings.TrimPrefix(h, "Bearer ")
		}
		if token == "" {
			if ck, err := c.Cookie("gotax_token"); err == nil {
				if u, err := url.QueryUnescape(ck); err == nil {
					token = u
				} else {
					token = ck
				}
			}
		}
		if token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		claims, err := auth.ParseAndValidateAccessToken(token)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}

// Page describes one server-rendered page. Load gathers the page's data.
type Page struct {
	Title   string
	NavPath string
	Load    func(c *gin.Context) (any, error)
}

// RegisterPages mounts one catch-all GET for /app/*. Converted pages (present
// in the template sets) are server-rendered; unconverted pages fall back to
// static files. Mutation endpoints are explicit POST routes (no conflict).
func (s *Server) RegisterPages(r *gin.Engine, pages map[string]Page, actions map[string]map[string]gin.HandlerFunc) {
	app := r.Group("/app")
	app.Use(PageAuthMiddleware())
	app.GET("/*filepath", func(c *gin.Context) {
		path := "/app" + c.Param("filepath")
		if path == "/app" || path == "/app/" {
			c.Redirect(http.StatusFound, "/app/dashboard.html")
			return
		}
		key := pageKey(path)
		if _, ok := s.sets[key]; ok {
			p, ok := pages[path]
			if !ok {
				http.NotFound(c.Writer, c.Request)
				return
			}
			shell := BuildShell(c, p.Title, p.NavPath)
			data, err := p.Load(c)
			if err != nil {
				log.Printf("load page %s: %v", path, err)
				shell.Title = "Lỗi"
				s.Render(c, key, shell, gin.H{"Error": err.Error()})
				return
			}
			s.Render(c, key, shell, data)
			return
		}
		http.ServeFile(c.Writer, c.Request, "web/app"+c.Param("filepath"))
	})
	for path, m := range actions {
		for action, fn := range m {
			r.POST(path+"/"+action, PageAuthMiddleware(), fn)
		}
	}
}

func pageKey(path string) string {
	base := path
	if i := strings.LastIndex(base, "/"); i >= 0 {
		base = base[i+1:]
	}
	base = strings.TrimSuffix(base, ".html")
	if strings.HasPrefix(path, "/payroll/") {
		return "payroll-" + base
	}
	return base
}
