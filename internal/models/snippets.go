package models

import (
	"database/sql"
	"errors"
	"time"
)

// Chapter 4.5: Designing a database model |
// Define a snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	title   string
	content string
	created time.Time
	expires time.Time
}

// Chapter 4.5: Designing a database model |
// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Chapter 4.5: Designing a database model |
// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Chapter 4.6: Executing SQL statements |
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	stmt := `INSERT INTO snippets(title, content, created, expires)
	VALUES(?, ?, NOW(), DATE_ADD(NOW(), INTERVAL ? DAY))`

	// Chapter 4.6: Executing SQL statements |
	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for placeholder parameters. This
	// method returns a sql.Result type, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Chapter 4.6: Executing SQL statements |
	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Chapter 4.6: Executing SQL statements |
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// Chapter 4.5: Designing a database model |
// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// Chapter 4.7: Single-record SQL queries |
	// Write the SQL statement we want to execute. Again,I've split it over three
	// lines for readability.
	stmt := `SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > NOW() AND id = ?`
	// Chapter 4.7: Single-record SQL queries |
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Chapter 4.7: Single-record SQL queries
	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Chapter 4.7: Single-record SQL queries
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.title, &s.content, &s.created, &s.expires)
	if err != nil {
		// Chapter 4.7: Single-record SQL queries |
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// Chapter 4.7: Single-record SQL queries
	// If everything went OK then return the Snippet object.
	return s, nil
}

// Chapter 4.5: Designing a database model |
// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
