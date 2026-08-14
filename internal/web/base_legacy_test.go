package web

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Legacy Alpine pages run inside the server shell. Alpine 3.15 auto-starts in
// a queueMicrotask after its <script> executes; without defer it boots before
// <body> exists and before auth-legacy/app-legacy register their alpine:init
// stores, so every legacy page loses its component state. The original static
// pages load Alpine with defer — the shell must too.
func TestBaseLegacyLoadsAlpineWithDefer(t *testing.T) {
	raw, err := os.ReadFile("web/templates/base-legacy.html")
	require.NoError(t, err)
	src := string(raw)
	i := strings.Index(src, "/assets/js/alpine.min.js")
	require.True(t, i > 0, "alpine.min.js script tag missing from base-legacy.html")
	line := src[strings.LastIndex(src[:i], "<script"):]
	require.Contains(t, line, "defer", "alpine.min.js must load with defer")
}
