// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"
)

type Track struct {
	Rank   int32
	Title  string
	Artist string
	Uri    string
	Date   time.Time
}
