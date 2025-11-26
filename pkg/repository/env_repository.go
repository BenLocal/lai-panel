package repository

import (
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
)

type EnvRepository struct {
	db *sqlx.DB
}

func NewEnvRepository() *EnvRepository {
	return &EnvRepository{db: database.GetDB()}
}

func (r *EnvRepository) Create(env *model.Env) error {
	query := `INSERT INTO env (key, value, scope, description, metadata) VALUES (:key, :value, :scope, :description, :metadata)`
	result, err := r.db.NamedExec(query, env)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	env.ID = id
	return nil
}

func (r *EnvRepository) Update(env *model.Env) error {
	query := `UPDATE env SET key = :key, value = :value, scope = :scope, description = :description WHERE id = :id`
	_, err := r.db.NamedExec(query, env)
	return err
}

func (r *EnvRepository) GetByKey(key string) (*model.Env, error) {
	query := `SELECT * FROM env WHERE key = ?`
	var env model.Env
	err := r.db.Get(&env, query, key)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *EnvRepository) GetScopes() ([]string, error) {
	query := `SELECT DISTINCT scope FROM env`
	var scopes []string
	err := r.db.Select(&scopes, query)
	if err != nil {
		return nil, err
	}
	return scopes, nil
}

func (r *EnvRepository) Delete(id int64) error {
	query := `DELETE FROM env WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EnvRepository) GetPage(scope string, page int, pageSize int) (int, []model.Env, error) {
	var total int
	query := `SELECT COUNT(*) FROM env`
	params := []interface{}{}
	if scope != "" {
		query += " WHERE scope = ?"
		params = append(params, scope)
	}

	err := r.db.Get(&total, query, params...)
	if err != nil {
		return 0, nil, err
	}

	if total == 0 {
		return 0, nil, nil
	}

	sql := `SELECT * FROM env`
	p := []interface{}{}
	if scope != "" {
		sql += " WHERE scope = ?"
		p = append(p, scope)
	}
	sql += " ORDER BY created_at DESC LIMIT ? OFFSET ?"

	var lst []model.Env
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	limit := pageSize
	offset := (page - 1) * pageSize
	p = append(p, limit, offset)
	err = r.db.Select(&lst, sql, p...)
	return total, lst, err
}
