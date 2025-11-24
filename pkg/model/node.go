package model

import (
	"time"

	"github.com/benlocal/lai-panel/pkg/crypto"
)

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
	DataPath    *string   `db:"data_path" json:"data_path"`
}

type NodeView struct {
	ID                 int64   `json:"id"`
	IsLocal            bool    `json:"is_local"`
	Name               string  `json:"name"`
	DisplayName        *string `json:"display_name"`
	Address            string  `json:"address"`
	Status             string  `json:"status"`
	SSHUser            string  `json:"ssh_user"`
	RequestSSHPassword *string `json:"ssh_password"`
	SSHPort            int     `json:"ssh_port"`
	AgentPort          int     `json:"agent_port"`
}

func (n *Node) ToView() *NodeView {
	return &NodeView{
		ID:                 n.ID,
		IsLocal:            n.IsLocal,
		Name:               n.Name,
		DisplayName:        n.DisplayName,
		Address:            n.Address,
		Status:             n.Status,
		SSHUser:            n.SSHUser,
		SSHPort:            n.SSHPort,
		AgentPort:          n.AgentPort,
		RequestSSHPassword: nil,
	}
}

func (n *Node) GetDecryptedSSHPassword() (string, error) {
	decrypted, err := crypto.Decrypt(n.SSHPassword)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func (v *NodeView) ToModel() (*Node, error) {
	var encryptedPassword string
	if v.RequestSSHPassword != nil && *v.RequestSSHPassword != "" {
		encrypted, err := crypto.Encrypt(*v.RequestSSHPassword)
		if err != nil {
			return nil, err
		}
		encryptedPassword = encrypted
	}

	return &Node{
		ID:          v.ID,
		IsLocal:     v.IsLocal,
		Name:        v.Name,
		DisplayName: v.DisplayName,
		Address:     v.Address,
		SSHUser:     v.SSHUser,
		SSHPassword: encryptedPassword,
		SSHPort:     v.SSHPort,
	}, nil
}
