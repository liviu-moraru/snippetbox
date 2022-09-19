package models

import (
	"database/sql"
	"time"
)

// User Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel Define a new UserModel type which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert We'll use the Insert method to add a new record to the "users" table.
func (m *UserModel) Insert(name string, email string, password string) error {
	return nil
}

// Authenticate We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

// Exists We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
