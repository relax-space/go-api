package controllers

const (
	// DefaultMaxResultCount In order to prevent the API from querying the data of the entire table, the default is 30
	DefaultMaxResultCount = 30
)

// SearchInput General query conditions
type SearchInput struct {
	Sortby         []string `query:"sortby"`
	Order          []string `query:"order"`
	SkipCount      int      `query:"skipCount"`
	MaxResultCount int      `query:"maxResultCount"`
	WithHasMore    bool     `query:"withHasMore"`
}
