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
	query := `INSERT INTO apps (name, display, docker_compose) 
	VALUES (:name, :display, :docker_compose)`

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
	query := `UPDATE apps SET name = :name, display = :display,
	 docker_compose = :docker_compose,
	  updated_at = CURRENT_TIMESTAMP WHERE id = :id`
	_, err := r.db.NamedExec(query, app)
	return err
}
