package ds

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FuncMap returns the shared template function map for all UXR applications.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// HTML safety
		"safe":     func(s string) template.HTML { return template.HTML(s) },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"safeAttr": func(s string) template.HTMLAttr { return template.HTMLAttr(s) },
		"safeURL":  func(s string) template.URL { return template.URL(s) },
		"safeJS":   func(s string) template.JS { return template.JS(s) },
		"safeCSS":  func(s string) template.CSS { return template.CSS(s) },

		// Dictionary — pass multiple values to sub-templates
		"dict": func(values ...any) map[string]any {
			if len(values)%2 != 0 {
				return nil
			}
			d := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					continue
				}
				d[key] = values[i+1]
			}
			return d
		},

		// Math
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},
		"seq": func(start, end int) []int {
			s := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},

		// Comparison
		"eq": func(a, b any) bool { return a == b },
		"ne": func(a, b any) bool { return a != b },
		"lt": func(a, b int) bool { return a < b },
		"le": func(a, b int) bool { return a <= b },
		"gt": func(a, b int) bool { return a > b },
		"ge": func(a, b int) bool { return a >= b },

		// Strings
		"truncate":  Truncate,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"title":     strings.Title, //nolint:staticcheck
		"trim":      strings.TrimSpace,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"split": strings.Split,
		"join":  strings.Join,

		// Time
		"formatTime":     FormatTime,
		"formatDate":     FormatDate,
		"formatDateTime": FormatDateTime,
		"timeAgo":        TimeAgo,
		"now":            time.Now,

		// Numbers
		"formatNumber": FormatNumber,

		// Conditional
		"default": func(defaultVal, val any) any {
			if isZeroValue(val) {
				return defaultVal
			}
			return val
		},
		"coalesce": func(vals ...any) any {
			for _, v := range vals {
				if !isZeroValue(v) {
					return v
				}
			}
			return nil
		},

		// Slices
		"first": func(v any) any {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Slice && rv.Len() > 0 {
				return rv.Index(0).Interface()
			}
			return nil
		},
		"last": func(v any) any {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Slice && rv.Len() > 0 {
				return rv.Index(rv.Len() - 1).Interface()
			}
			return nil
		},
		"length": func(v any) int {
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Slice, reflect.Map, reflect.String:
				return rv.Len()
			default:
				return 0
			}
		},
	}
}

// MergeFuncMap creates a new FuncMap from base, overriding/adding entries from extra.
func MergeFuncMap(base, extra template.FuncMap) template.FuncMap {
	merged := make(template.FuncMap, len(base)+len(extra))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}

// toTime coerces any to time.Time. Accepts time.Time, RFC3339/date string, or Unix int64.
func toTime(v any) (time.Time, bool) {
	switch t := v.(type) {
	case time.Time:
		if t.IsZero() {
			return time.Time{}, false
		}
		return t, true
	case string:
		if t == "" {
			return time.Time{}, false
		}
		for _, layout := range []string{time.RFC3339, "2006-01-02"} {
			if parsed, err := time.Parse(layout, t); err == nil {
				return parsed, true
			}
		}
		return time.Time{}, false
	case int64:
		if t == 0 {
			return time.Time{}, false
		}
		return time.Unix(t, 0), true
	default:
		return time.Time{}, false
	}
}

// formatTimeValue formats any time-like value with the given layout.
func formatTimeValue(v any, layout string) string {
	t, ok := toTime(v)
	if !ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
		return ""
	}
	return t.Format(layout)
}

// FormatTime formats a time value for display (date + time).
func FormatTime(v any) string { return formatTimeValue(v, "Jan 2, 2006 3:04 PM") }

// FormatDate formats a time value as date only.
func FormatDate(v any) string { return formatTimeValue(v, "Jan 02, 2006") }

// FormatDateTime formats a time value as date and time (24h).
func FormatDateTime(v any) string { return formatTimeValue(v, "Jan 02, 2006 15:04") }

// TimeAgo returns a human-readable relative time string.
func TimeAgo(v any) string {
	t, ok := toTime(v)
	if !ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
		return ""
	}

	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	case d < 30*24*time.Hour:
		w := int(d.Hours() / 24 / 7)
		if w == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", w)
	case d < 365*24*time.Hour:
		mo := int(d.Hours() / 24 / 30)
		if mo == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", mo)
	default:
		y := int(d.Hours() / 24 / 365)
		if y == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", y)
	}
}

// FormatNumber formats an int64 with comma separators.
func FormatNumber(n int64) string {
	s := strconv.FormatInt(n, 10)
	if n < 0 {
		return "-" + insertCommas(s[1:])
	}
	return insertCommas(s)
}

func insertCommas(s string) string {
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	start := len(s) % 3
	if start > 0 {
		b.WriteString(s[:start])
	}
	for i := start; i < len(s); i += 3 {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

// isZeroValue returns true if v is nil or a zero value for its type.
func isZeroValue(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Slice, reflect.Map:
		return rv.Len() == 0
	default:
		return false
	}
}
