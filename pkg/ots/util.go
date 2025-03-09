package ots

import "github.com/rs/zerolog/log"

func MapMustGetString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		} else {
			log.Warn().Str("key", key).Interface("map", m).Interface("value", v).Msg("MapMustGetString: type assertion failed")
		}
	}
	return ""
}
