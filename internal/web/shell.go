package web

import (
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// Shell carries the shared layout data (sidebar, topbar).
type Shell struct {
	Title       string
	NavPath     string
	Username    string
	RoleLabel   string
	AvatarInit  string
	CompanyName string
}

// BuildShell derives layout data from the gin context (set by PageAuth).
func BuildShell(c *gin.Context, title, navPath string) Shell {
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	role := c.GetString("role")
	return Shell{
		Title:       title,
		NavPath:     navPath,
		Username:    username,
		RoleLabel:   roleLabel(role),
		AvatarInit:  avatarInit(username),
		CompanyName: "GoTax",
	}
}

func avatarInit(username string) string {
	if username == "" {
		return "U"
	}
	r := []rune(username)[0]
	return string(unicode.ToUpper(r))
}

func roleLabel(role string) string {
	switch role {
	case "admin":
		return "Admin — Quản trị viên"
	case "chief_accountant":
		return "Kế toán trưởng"
	case "accountant":
		return "Kế toán viên"
	case "viewer":
		return "Viewer — Chỉ xem"
	default:
		return strings.ToUpper(role)
	}
}
