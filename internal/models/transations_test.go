package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestExampleTransaction(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO snippets").WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectExec("UPDATE snippets").WithArgs("New_Title", 10).WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectCommit()

	// Act
	m := &SnippetModel{DB: db}
	var id int
	if id, err = m.ExampleTransaction(); err != nil {
		t.Errorf("error was not expected while executing ExampleTransaction: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if id != 10 {
		t.Errorf("there returned value of id is not correct: %d", id)
	}
}
