package repository

import (
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
)

type NodeRepository struct {
	db *sqlx.DB
}

func NewNodeRepository() *NodeRepository {
	return &NodeRepository{db: database.GetDB()}
}

func (r *NodeRepository) Create(node *model.Node) error {
	query := `INSERT INTO nodes (name, address, ssh_port,
	 ssh_user, ssh_password, agent_port, status, is_local) 
	          VALUES (:name, :address, :ssh_port, :ssh_user, 
			  :ssh_password, :agent_port, :status, :is_local) RETURNING id`

	result, err := r.db.NamedExec(query, node)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	node.ID = id
	return nil
}

func (r *NodeRepository) GetByID(id int64) (*model.Node, error) {
	var node model.Node
	err := r.db.Get(&node, "SELECT * FROM nodes WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NodeRepository) GetByNodeName(name string) (*model.Node, error) {
	var node model.Node
	err := r.db.Get(&node, "SELECT * FROM nodes WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NodeRepository) Update(node *model.Node) error {
	query := `UPDATE nodes SET name = :name,
	 display_name = :display_name,
	 address = :address, 
	 ssh_port = :ssh_port,
	 ssh_user = :ssh_user,
	 ssh_password = :ssh_password,
	 updated_at = CURRENT_TIMESTAMP 
	 WHERE id = :id`

	_, err := r.db.NamedExec(query, node)
	return err
}

func (r *NodeRepository) UpdateRegistry(node *model.Node) error {
	query := `UPDATE nodes SET status = :status,
	 updated_at = CURRENT_TIMESTAMP 
	 WHERE id = :id`

	_, err := r.db.NamedExec(query, node)
	return err
}

func (r *NodeRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM nodes WHERE id = ?", id)
	return err
}

func (r *NodeRepository) List() ([]model.Node, error) {
	var nodes []model.Node
	err := r.db.Select(&nodes, "SELECT * FROM nodes ORDER BY created_at DESC")
	return nodes, err
}
