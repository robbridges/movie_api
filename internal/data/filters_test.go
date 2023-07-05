package data

import (
	"movie_api/internal/validator"
	"testing"
)

func TestValidateFilters(t *testing.T) {
	t.Run("Happy path, all correct", func(t *testing.T) {
		filters := Filters{
			Page:         2,
			PageSize:     3,
			Sort:         "id",
			SortSafeList: []string{"id"},
		}
		validator := validator.New()
		ValidateFilters(validator, filters)
		if !validator.Valid() {
			t.Errorf("Validator contains errors when it should not")
		}
	})
}

func TestSortColumnSafeValue(t *testing.T) {
	filters := Filters{
		Sort:         "safe",
		SortSafeList: []string{"safe"},
	}
	// no panic should be reached, this is safe
	defer func() {
		if r := recover(); r != nil {
			t.Error("Expected no panic, but a panic occurred")
		}
	}()

	filters.SortColumn()
}

func TestSortColumnUnsafeValue(t *testing.T) {
	filters := Filters{
		Sort:         "unsafe",
		SortSafeList: []string{"safe"},
	}

	// we expect this panic to hit.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic, but no panic occurred")
		}
	}()

	filters.SortColumn()
}

func TestFilters_SortDirection(t *testing.T) {
	t.Run("Should return ASC", func(t *testing.T) {
		filter := Filters{
			Sort:         "runtime",
			SortSafeList: []string{"runtime, -runtime"},
		}
		want := "ASC"
		got := filter.SortDirection()
		if want != got {
			t.Errorf("Test sort Direction: got: %s, want %s", got, want)
		}
	})
	t.Run("Should return DESC", func(t *testing.T) {
		filter := Filters{
			Sort:         "-runtime",
			SortSafeList: []string{"runtime, -runtime"},
		}
		want := "DESC"
		got := filter.SortDirection()
		if want != got {
			t.Errorf("Test sort Direction: got: %s, want %s", got, want)
		}
	})
}

func TestValidateFiltersSadPaths(t *testing.T) {
	tests := []struct {
		name       string
		errorCount int
		filters    Filters
	}{
		{
			name:       "Sad path, page too small",
			errorCount: 1,
			filters: Filters{
				Page:         -2,
				PageSize:     3,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
		},
		{
			name:       "Sad path, pagesize too small",
			errorCount: 1,
			filters: Filters{
				Page:         2,
				PageSize:     -3,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
		},
		{
			name:       "Sad Path, bad sort param",
			errorCount: 1,
			filters: Filters{
				Page:         2,
				PageSize:     3,
				Sort:         "id",
				SortSafeList: []string{"accepted"},
			},
		},
		{
			name:       "Sad path, page too big",
			errorCount: 1,
			filters: Filters{
				Page:         100_000_000,
				PageSize:     3,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
		},
		{
			name:       "Sad path, pagesize too big",
			errorCount: 1,
			filters: Filters{
				Page:         1,
				PageSize:     3000,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
		},
		{
			name:       "Sad path, 2 params wrong",
			errorCount: 2,
			filters: Filters{
				Page:         -1,
				PageSize:     3000,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
		},
		{
			name:       "Sad path, all params wrong",
			errorCount: 3,
			filters: Filters{
				Page:         -1,
				PageSize:     3000,
				Sort:         "id",
				SortSafeList: []string{"accepted"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			validator := validator.New()
			ValidateFilters(validator, tc.filters)
			if validator.Valid() {
				t.Errorf("Validator Should contain an error and it did not")
			}
			want := tc.errorCount
			got := len(validator.Errors)
			if len(validator.Errors) != tc.errorCount {
				t.Errorf("Wrong amount of errors got %d want %d", got, want)
			}
		})
	}
}
