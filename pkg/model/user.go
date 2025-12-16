package model

import "time"

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
