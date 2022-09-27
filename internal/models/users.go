package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
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
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Use the Exec() method to insert the user details and hashed password
	// into the users table.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email string, password string) (User, error) {
	stmt := `SELECT id, name, hashed_password FROM users
				WHERE email = ?`
	u := User{}
	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, ErrInvalidCredentials
		}
		return u, err
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return u, ErrInvalidCredentials
		} else {
			return u, err
		}

	}
	return u, nil
}

// Exists We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
