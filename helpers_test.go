package ds

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		s      string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"hi", 2, "hi"},
		{"hello", 3, "hel"},
		{"hello", 5, "hello"},
	}
	for _, tt := range tests {
		got := Truncate(tt.s, tt.maxLen)
		if got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
		}
	}
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		current, total int
		wantLen        int
	}{
		{1, 3, 3},
		{1, 7, 7},
		{1, 10, 6},  // pages 1-5 + 10
		{5, 10, 5},  // 1, 4, 5, 6, 10
		{9, 10, 6},  // 1, 6-10
	}
	for _, tt := range tests {
		p := NewPagination(tt.current, tt.total)
		if len(p.PageNumbers) != tt.wantLen {
			t.Errorf("NewPagination(%d, %d) got %d pages %v, want %d",
				tt.current, tt.total, len(p.PageNumbers), p.PageNumbers, tt.wantLen)
		}
		if p.CurrentPage != tt.current {
			t.Errorf("CurrentPage = %d, want %d", p.CurrentPage, tt.current)
		}
	}
}

func TestBoolState(t *testing.T) {
	if got := BoolState(true); got != "pass" {
		t.Errorf("BoolState(true) = %q", got)
	}
	if got := BoolState(false); got != "fail" {
		t.Errorf("BoolState(false) = %q", got)
	}
	if got := BoolState(false, "warn"); got != "warn" {
		t.Errorf("BoolState(false, warn) = %q", got)
	}
}

func TestBoolYesNo(t *testing.T) {
	if got := BoolYesNo(true); got != "Yes" {
		t.Errorf("BoolYesNo(true) = %q", got)
	}
	if got := BoolYesNo(false); got != "No" {
		t.Errorf("BoolYesNo(false) = %q", got)
	}
}
