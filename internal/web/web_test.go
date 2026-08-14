package web

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestMain runs every web test from the repo root so template/static paths
// (web/templates, web/app) resolve the same way as a production `go run .`.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	if err := os.Chdir("../.."); err != nil {
		panic("chdir repo root: " + err.Error())
	}
	os.Exit(m.Run())
}
