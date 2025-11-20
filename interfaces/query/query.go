package query

// SortDirection defines the direction of sorting.
type SortDirection string

const (
	ASC  SortDirection = "ASC"  // Ascending order
	DESC SortDirection = "DESC" // Descending order
)

// SortOrder defines a single sorting criterion.
type SortOrder struct {
	Field     string        // Field to sort by
	Direction SortDirection // Sorting direction (ASC or DESC)
}

// PageRequest defines parameters for pagination.
type PageRequest struct {
	Page     int         // Current page number (1-based)
	PageSize int         // Number of items per page
	SortBy   []SortOrder // Sorting criteria
}

// PageResponse holds paginated results.
type PageResponse[T any] struct {
	Content       []T   // The actual content for the current page
	TotalElements int64 // Total number of elements across all pages
	TotalPages    int   // Total number of pages
	Page          int   // Current page number
	PageSize      int   // Number of elements per page
}
