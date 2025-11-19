package repository

import (
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
)

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{db: database.GetDB()}
}

func (r *ServiceRepository) Create(service *model.Service) error {
	query := `INSERT INTO services (name, app_id, node_id, status, metadata) 
	VALUES (:name, :app_id, :node_id, :status, :metadata)`
	_, err := r.db.NamedExec(query, service)
	return err
}

func (r *ServiceRepository) GetByID(id int64) (*model.Service, error) {
	query := `SELECT * FROM services WHERE id = ?`
	var service model.Service
	err := r.db.Get(&service, query, id)
	return &service, err
}

func (r *ServiceRepository) Update(service *model.Service) error {
	query := `UPDATE services SET name = :name, app_id = :app_id, 
	node_id = :node_id,
	metadata = :metadata,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = :id`
	_, err := r.db.NamedExec(query, service)
	return err
}

func (r *ServiceRepository) Delete(id int64) error {
	query := `DELETE FROM services WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ServiceRepository) GetPage(page, pageSize int) (int, []*model.Service, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM services")
	if err != nil {
		return 0, nil, err
	}

	if total == 0 {
		return 0, nil, nil
	}

	query := `SELECT * FROM services ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var services []*model.Service
	limit := pageSize
	offset := (page - 1) * pageSize
	err = r.db.Select(&services, query, limit, offset)
	return total, services, err
}
