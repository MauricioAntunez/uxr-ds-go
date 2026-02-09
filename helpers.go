package ds

import "html/template"

// Metric represents a single metric display.
type Metric struct {
	ID    string
	Value string
	Label string
}

// MetricsGrid represents a grid of metrics.
type MetricsGrid struct {
	ID      string
	Metrics []Metric
}

// Column represents a table column.
type Column struct {
	Label string
	Width string
}

// DataTable represents a data table component.
type DataTable struct {
	ID      string
	Columns []Column
	Rows    template.HTML
}

// Tab represents a single tab.
type Tab struct {
	ID      string
	Label   string
	Content template.HTML
}

// Tabs represents a tab container.
type Tabs struct {
	Tabs []Tab
}

// PageHeader represents a page header component.
type PageHeader struct {
	Title    string
	Subtitle string
	Actions  template.HTML
}

// FormField represents a form field component.
type FormField struct {
	ID          string
	Name        string
	Label       string
	Type        string // text, email, password, select, textarea
	Placeholder string
	Required    bool
	Options     []SelectOption
}

// SelectOption represents a select dropdown option.
type SelectOption struct {
	Value string
	Label string
}

// Card represents a card component.
// Supports both .Content and .Body (alias), both .HeaderAction and .HeaderActions.
type Card struct {
	Header        string
	HeaderAction  template.HTML
	HeaderActions template.HTML
	Body          template.HTML
	Content       template.HTML
	Elevated      bool
	Class         string
}

// StatusBadge represents a status badge component.
type StatusBadge struct {
	Status string
	Label  string
}

// EmptyState represents an empty state component.
type EmptyState struct {
	Message   string
	Action    string
	ActionURL string
}

// Pagination represents pagination component data.
type Pagination struct {
	CurrentPage int
	TotalPages  int
	PageNumbers []int
}

// ConfirmDialog represents a native <dialog> confirmation.
type ConfirmDialog struct {
	ID           string
	Title        string
	Message      string
	ConfirmLabel string
}

// PopoverMenu represents a Popover API-based dropdown menu.
type PopoverMenu struct {
	ID      string
	Label   string
	Content template.HTML
}

// NewPagination creates pagination data with smart page number ranges.
func NewPagination(currentPage, totalPages int) Pagination {
	pageNumbers := make([]int, 0, totalPages)

	if totalPages <= 7 {
		for i := 1; i <= totalPages; i++ {
			pageNumbers = append(pageNumbers, i)
		}
	} else {
		if currentPage <= 4 {
			for i := 1; i <= 5; i++ {
				pageNumbers = append(pageNumbers, i)
			}
			pageNumbers = append(pageNumbers, totalPages)
		} else if currentPage >= totalPages-3 {
			pageNumbers = append(pageNumbers, 1)
			for i := totalPages - 4; i <= totalPages; i++ {
				pageNumbers = append(pageNumbers, i)
			}
		} else {
			pageNumbers = append(pageNumbers, 1)
			for i := currentPage - 1; i <= currentPage+1; i++ {
				pageNumbers = append(pageNumbers, i)
			}
			pageNumbers = append(pageNumbers, totalPages)
		}
	}

	return Pagination{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		PageNumbers: pageNumbers,
	}
}

// Truncate truncates a string to maxLen characters with ellipsis.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// BoolState returns "pass" if ok, otherwise the fallback or "fail".
func BoolState(ok bool, fallback ...string) string {
	if ok {
		return "pass"
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return "fail"
}

// BoolYesNo returns "Yes" if ok, otherwise "No".
func BoolYesNo(ok bool) string {
	if ok {
		return "Yes"
	}
	return "No"
}
