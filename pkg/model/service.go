package model

import "time"

type Service struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	AppID     int64     `db:"app_id" json:"app_id"`
	NodeID    int64     `db:"node_id" json:"node_id"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Metadata  *string   `db:"metadata" json:"metadata"`
}
