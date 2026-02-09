package ds

import (
	"html/template"
	"testing"
	"time"
)

func TestFuncMap(t *testing.T) {
	fm := FuncMap()
	if fm == nil {
		t.Fatal("FuncMap returned nil")
	}

	// Spot-check key functions exist
	for _, name := range []string{
		"safe", "safeHTML", "dict", "add", "sub", "mul", "div", "mod", "seq",
		"eq", "ne", "lt", "le", "gt", "ge",
		"truncate", "lower", "upper", "title", "trim", "contains",
		"hasPrefix", "hasSuffix", "replace", "split", "join",
		"formatTime", "formatDate", "formatDateTime", "timeAgo", "now",
		"formatNumber", "default", "coalesce", "first", "last", "length",
	} {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap missing function: %s", name)
		}
	}
}

func TestMergeFuncMap(t *testing.T) {
	base := template.FuncMap{"a": func() string { return "base" }}
	extra := template.FuncMap{"b": func() string { return "extra" }, "a": func() string { return "overridden" }}

	merged := MergeFuncMap(base, extra)
	if _, ok := merged["a"]; !ok {
		t.Error("merged missing 'a'")
	}
	if _, ok := merged["b"]; !ok {
		t.Error("merged missing 'b'")
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"zero time", time.Time{}, ""},
		{"empty string", "", ""},
		{"int64 zero", int64(0), ""},
		{"nil", nil, ""},
		{"rfc3339 string", "2024-06-15T10:30:00Z", "Jun 15, 2024 10:30 AM"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTime(tt.input)
			if got != tt.want {
				t.Errorf("FormatTime(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"zero time", time.Time{}, ""},
		{"date string", "2024-06-15", "Jun 15, 2024"},
		{"rfc3339 string", "2024-06-15T10:30:00Z", "Jun 15, 2024"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDate(tt.input)
			if got != tt.want {
				t.Errorf("FormatDate(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	tm := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	got := FormatDateTime(tm)
	want := "Jun 15, 2024 14:30"
	if got != want {
		t.Errorf("FormatDateTime = %q, want %q", got, want)
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1,000"},
		{1234567, "1,234,567"},
		{-1234, "-1,234"},
	}
	for _, tt := range tests {
		got := FormatNumber(tt.input)
		if got != tt.want {
			t.Errorf("FormatNumber(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTimeAgo(t *testing.T) {
	if got := TimeAgo(""); got != "" {
		t.Errorf("TimeAgo empty string = %q", got)
	}
	if got := TimeAgo(nil); got != "" {
		t.Errorf("TimeAgo nil = %q", got)
	}
	if got := TimeAgo(time.Now().Add(-30 * time.Second)); got != "just now" {
		t.Errorf("TimeAgo 30s ago = %q, want 'just now'", got)
	}
	if got := TimeAgo(time.Now().Add(-5 * time.Minute)); got != "5 minutes ago" {
		t.Errorf("TimeAgo 5m ago = %q", got)
	}
}
