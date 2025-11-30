package i18n

import (
	"os"
	"strings"
)

var (
	CurrentLang = "zh"
)

func init() {
	lang := os.Getenv("AI_OPS_LANG")
	if lang != "" {
		CurrentLang = strings.ToLower(lang)
	}
}

func T(key string) string {
	if msgs, ok := messages[CurrentLang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to default (zh) if not found in current lang
	if CurrentLang != "zh" {
		if msgs, ok := messages["zh"]; ok {
			if msg, ok := msgs[key]; ok {
				return msg
			}
		}
	}
	return key
}
