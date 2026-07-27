package i18n

import (
	"embed"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed *.json
var fs embed.FS

type Localizer struct {
	bundle *i18n.Bundle
}

func New() (*Localizer, error) {
	b := i18n.NewBundle(language.English)
	entries, err := fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("read i18n dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if _, err := b.LoadMessageFileFS(fs, e.Name()); err != nil {
			return nil, fmt.Errorf("load %s: %w", e.Name(), err)
		}
	}
	return &Localizer{bundle: b}, nil
}

func MustNew() *Localizer {
	l, err := New()
	if err != nil {
		panic(err)
	}
	return l
}

func (l *Localizer) ForLocale(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(l.bundle, lang)
}

func Localize(localizer *i18n.Localizer, msgID string, args ...map[string]any) string {
	if localizer == nil {
		return msgID
	}
	data := map[string]any{}
	if len(args) > 0 {
		data = args[0]
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: data,
	})
	if err != nil {
		return msgID
	}
	return msg
}
