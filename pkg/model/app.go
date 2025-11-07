package model

import "time"

type App struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Display       *string   `db:"display" json:"display"`
	DockerCompose *string   `db:"docker_compose" json:"docker_compose"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	Metadata      *string   `db:"metadata" json:"metadata"`
}
