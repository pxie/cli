package api

import (
	"fmt"
	"generic"
)

func stringOrNull(s interface{}) string {
	switch s := s.(type) {
	case string:
		if s == "" {
			return "null"
		}
		return s
	default:
		return fmt.Sprintf(`%s`, s)
	}
}

func mapToJsonValues(params generic.Map) (vals []string) {
	generic.Each(params, func(key, val interface{}) {
		switch val := val.(type) {
		case string:
			if val != "null" {
				val = fmt.Sprintf(`"%s"`, val)
			}
			vals = append(vals, fmt.Sprintf(`"%s":%s`, key, val))
		case int:
			vals = append(vals, fmt.Sprintf(`"%s":%d`, key, val))
		case uint64:
			vals = append(vals, fmt.Sprintf(`"%s":%d`, key, val))
		default:
			vals = append(vals, fmt.Sprintf(`"%s":%s`, key, val))
		}
	})
	return
}
