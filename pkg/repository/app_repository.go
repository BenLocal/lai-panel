package repository

import (
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
)

type AppRepository struct {
	db *sqlx.DB
}

func NewAppRepository() *AppRepository {
	return &AppRepository{db: database.GetDB()}
}

func (r *AppRepository) Create(app *model.App) error {
	query := `INSERT INTO apps (name, display, version, icon, docker_compose, metadata, qa, description) 
	VALUES (:name, :display, :version, :icon, :docker_compose, :metadata, :qa, :description)`

	result, err := r.db.NamedExec(query, app)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	app.ID = id
	return nil
}

func (r *AppRepository) GetByID(id int64) (*model.App, error) {
	query := `SELECT * FROM apps WHERE id = ?`
	var app model.App
	err := r.db.Get(&app, query, id)
	return &app, err
}

func (r *AppRepository) Update(app *model.App) error {
	query := `UPDATE apps SET name = :name, 
		display = :display,
		version = :version,
		icon = :icon,
		docker_compose = :docker_compose,
		metadata = :metadata,
		qa = :qa,
		description = :description,
	  updated_at = CURRENT_TIMESTAMP WHERE id = :id`
	_, err := r.db.NamedExec(query, app)
	return err
}

func (r *AppRepository) List() ([]model.App, error) {
	query := `SELECT * FROM apps ORDER BY created_at DESC`
	var apps []model.App
	err := r.db.Select(&apps, query)
	return apps, err
}

func (r *AppRepository) Delete(id int64) error {
	query := `DELETE FROM apps WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
