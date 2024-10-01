package tests

import (
	"database/sql"
	"testing"
)

// TestOrderBy tests various ORDER BY scenarios
func TestOrderBy(t *testing.T) {
	// Connect to the database
	db, err := sql.Open("GoSql", "memory")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set up test data
	setupTestData(t, db)

	testCases := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "Order by id descending",
			query:    "SELECT name FROM users ORDER BY id DESC",
			expected: []string{"Charlie", "Alice", "Bob"},
		},
		{
			name:     "Order by number ascending",
			query:    "SELECT name FROM users ORDER BY 1 DESC",
			expected: []string{"Charlie", "Bob", "Alice"},
		},
		{
			name:     "Order sum of id and name by number ascending",
			query:    "SELECT '' + id + name FROM users ORDER BY 1 DESC",
			expected: []string{"2Charlie", "1Alice", "0Bob"},
		},
		{
			name:     "Order by fieldname descending",
			query:    "SELECT name FROM users ORDER BY name DESC",
			expected: []string{"Charlie", "Bob", "Alice"},
		},
		{
			name:     "Order by alias",
			query:    "SELECT name AS user_name FROM users ORDER BY user_name",
			expected: []string{"Alice", "Bob", "Charlie"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rows, err := db.Query(tc.query)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			defer rows.Close()

			var results []string
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					t.Fatalf("Failed to scan row: %v", err)
				}
				results = append(results, name)
			}

			if err := rows.Err(); err != nil {
				t.Fatalf("Error iterating over rows: %v", err)
			}

			if !compareSlices(results, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, results)
			}
		})
	}
}

func setupTestData(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT
        )`)

	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}
	_, err = db.Exec(
		"INSERT INTO users (name) VALUES ('Bob'), ('Alice'), ('Charlie')",
	)

	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}
}

func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}