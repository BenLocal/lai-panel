package repository

import (
	"github.com/benlocal/lai-panel/pkg/crypto"
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

func (r *UserRepository) Create(user *model.User) (int64, error) {
	password, err := crypto.Encrypt(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = password
	query := `INSERT INTO users (username, email, role, name, password) VALUES
	(:username, :email, :role, :name, :password)`
	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `SELECT * FROM users WHERE username = ?`
	var user model.User
	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetById(id int64) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	var user model.User
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdatePassword(id int64, password string) error {
	p, err := crypto.Encrypt(password)
	if err != nil {
		return err
	}
	query := `UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = r.db.Exec(query, p, id)
	return err
}
