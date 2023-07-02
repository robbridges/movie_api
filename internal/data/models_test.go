package data

import (
	"database/sql"
	"testing"
)

func TestNewModels(t *testing.T) {
	db := &sql.DB{} // Mocked database connection

	models := NewModels(db)

	// Check that the Movies field is initialized correctly
	if models.Movies.DB != db {
		t.Errorf("Expected Movies.DB to be %v, got %v", db, models.Movies.DB)
	}
}
