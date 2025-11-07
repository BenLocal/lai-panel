package model

import "time"

type Node struct {
	ID          int64     `db:"id" json:"id"`
	IsLocal     bool      `db:"is_local" json:"is_local"`
	Name        string    `db:"name" json:"name"`
	DisplayName *string   `db:"display_name" json:"display_name"`
	Address     string    `db:"address" json:"address"`
	SSHPort     int       `db:"ssh_port" json:"ssh_port"`
	AgentPort   int       `db:"agent_port" json:"agent_port"`
	SSHUser     string    `db:"ssh_user" json:"ssh_user"`
	SSHPassword string    `db:"ssh_password" json:"ssh_password"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Metadata    *string   `db:"metadata" json:"metadata"`
}
