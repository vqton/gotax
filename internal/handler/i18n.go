package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	gotaxi18n "gotax/internal/i18n"
)

const localizerKey = "localizer"

func I18nMiddleware(l *gotaxi18n.Localizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		c.Set(localizerKey, l.ForLocale(lang))
		c.Next()
	}
}

func T(c *gin.Context, msgID string, args ...map[string]any) string {
	v, exists := c.Get(localizerKey)
	if !exists {
		return msgID
	}
	localizer, ok := v.(*i18n.Localizer)
	if !ok {
		return msgID
	}
	return gotaxi18n.Localize(localizer, msgID, args...)
}
