package i18n

import (
	"ai-ops-agent/pkg/env"
	"strings"
)

var (
	CurrentLang = "en"
)

func init() {
	CurrentLang = strings.ToLower(env.Get("AI_OPS_LANG", "en"))
}

func T(key string) string {
	if msgs, ok := messages[CurrentLang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to default (en) if not found in current lang
	if CurrentLang != "en" {
		if msgs, ok := messages["en"]; ok {
			if msg, ok := msgs[key]; ok {
				return msg
			}
		}
	}
	return key
}
