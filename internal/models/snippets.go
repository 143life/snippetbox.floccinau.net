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

// *Chapter 4.9: Transactions and other details |
// We need somewhere to store the prepared statement for the lifetime of our
// web application. A neat way is to embed in the model alongside the connection
// pool.
// Chapter 4.5: Designing a database model |
// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB         *sql.DB
	InsertStmt *sql.Stmt
	GetStmt    *sql.Stmt
	LatestStmt *sql.Stmt
}

// *Chapter 4.9: Transactions and other details |
// Create a constructor for the model, in which we set up the prepared
// statement.
func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {
	// *Chapter 4.9: Transactions and other details |
	// Use the Prepare method to create a new prepared statement for the
	// current connection pool. This returns a sql.Stmt object which represents
	// the prepared statement
	var insertStmt, getStmt, latestStmt *sql.Stmt
	var err error
	insertStmt, err = db.Prepare(
		`INSERT INTO snippets(title, content, created, expires)
		VALUES(?, ?, NOW(), DATE_ADD(NOW(), INTERVAL ? DAY))`,
	)
	if err != nil {
		return nil, err
	}

	getStmt, err = db.Prepare(
		`SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > NOW() AND id = ?`,
	)
	if err != nil {
		return nil, err
	}

	latestStmt, err = db.Prepare(
		`SELECT id, title, content, created, expires
		FROM snippets
		ORDER BY id DESC LIMIT 10`,
	)
	if err != nil {
		return nil, err
	}

	// *Chapter 4.9: Transactions and other details |
	// Store it in our SnippetModel object, alongside the connection pool.
	return &SnippetModel{
		DB:         db,
		InsertStmt: insertStmt,
		GetStmt:    getStmt,
		LatestStmt: latestStmt,
	}, nil
}

// Chapter 4.5: Designing a database model |
// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Chapter 4.6: Executing SQL statements |
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	// stmt := `INSERT INTO snippets(title, content, created, expires)
	// VALUES(?, ?, NOW(), DATE_ADD(NOW(), INTERVAL ? DAY))`

	// Chapter 4.6: Executing SQL statements |
	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for placeholder parameters. This
	// method returns a sql.Result type, which contains some basic
	// information about what happened when the statement was executed.
	// result, err := m.DB.Exec(stmt, title, content, expires)
	// if err != nil {
	// 	return 0, err
	// }

	// *Chapter 4.9: Transactions and other details |
	// Notice how we call Exec directly against the prepared statement, rather
	// than against the connection pool? Prepared statements also support the
	// Query and QueryRow methods
	result, err := m.InsertStmt.Exec(title, content, expires)
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
	// stmt := `SELECT id, title, content, created, expires
	// FROM snippets
	// WHERE expires > NOW() AND id = ?`
	// Chapter 4.7: Single-record SQL queries |
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	// row := m.DB.QueryRow(stmt, id)

	// *Chapter 4.9: Transactions and other details |
	row := m.GetStmt.QueryRow(id)

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
	// Chapter 4.8: Multiple-record SQL queries |
	//  Write the SQL statement we want to execute
	// stmt := `SELECT id, title, content, created, expires
	// FROM snippets
	// WHERE expires > NOW()
	// ORDER BY id DESC LIMIT 10`

	// Chapter 4.8: Multiple-record SQL queries |
	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	// rows, err := m.DB.Query(stmt)
	// if err != nil {
	// 	return nil, err
	// }

	// *Chapter 4.9: Transactions and other details |
	rows, err := m.LatestStmt.Query()
	if err != nil {
		return nil, err
	}

	// Chapter 4.8: Multiple-record SQL queries |
	// We defer rows.Close() to ensure the SQL.rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Chapter 4.8: Multiple-record SQL queries |
	// Initialize an empty slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// Chapter 4.8: Multiple-record SQL queries |
	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Chapter 4.8: Multiple-record SQL queries |
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}
		// Chapter 4.8: Multiple-record SQL queries |
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.title, &s.content, &s.created, &s.expires)
		if err != nil {
			return nil, err
		}

		// Chapter 4.8: Multiple-record SQL queries |
		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// Chapter 4.8: Multiple-record SQL queries |
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error thet was encountered during the iteration. it's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Chapter 4.8: Multiple-record SQL queries
	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
