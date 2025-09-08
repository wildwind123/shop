package db

import "time"

type Product struct {
	ID          int       `db:"id" fieldtag:"db"`
	Name        string    `db:"name" fieldtag:"db"`
	Description string    `db:"description" fieldtag:"db"`
	CreatedAt   time.Time `db:"created_at" fieldtag:"db"`
	UpdatedAt   time.Time `db:"updated_at" fieldtag:"db"`
}
