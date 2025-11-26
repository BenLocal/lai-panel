package model

import "time"

type Env struct {
	ID          int64     `db:"id" json:"id"`
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	Scope       string    `db:"scope" json:"scope"`
	Description string    `db:"description" json:"description"`
	Metadata    string    `db:"metadata" json:"metadata"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
