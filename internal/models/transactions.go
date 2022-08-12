package models

func (m *SnippetModel) ExampleTransaction() (int, error) {
	// Calling the Begin() method on the connection pool creates a new sql.Tx
	// object, which represents the in-progress database transaction.
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Defer a call to tx.Rollback() to ensure it is always called before the
	// function returns. If the transaction succeeds it will be already be
	// committed by the time tx.Rollback() is called, making tx.Rollback() a
	// no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
	// the changes before the function returns.
	defer tx.Rollback()

	// Call Exec() on the transaction, passing in your statement and any
	// parameters. It's important to notice that tx.Exec() is called on the
	// transaction object just created, NOT the connection pool. Although we're
	// using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
	// exactly the same way.
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := tx.Exec(stmt, "Title", "Content", 1)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	// Carry out another transaction in exactly the same way.
	stmt = "UPDATE snippets SET title=? WHERE id=?"
	_, err = tx.Exec(stmt, "New_Title", id)
	if err != nil {
		return 0, err
	}

	// If there are no errors, the statements in the transaction can be committed
	// to the database with the tx.Commit() method.
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return int(id), err
}
