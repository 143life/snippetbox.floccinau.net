package models

import (
	"errors"
)

// Chapter 4.7: Single-record SQL queries
var ErrNoRecord = errors.New("models: no matching record found")
