package database

import (
	"github.com/benlocal/lai-panel/pkg/crypto"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
)

const (
	AdminUsername        = "admin"
	AdminEmail           = "admin@example.com"
	AdminRole            = "admin"
	AdminName            = "Admin"
	DefaultAdminPassword = "admin"
)

func ensureAdminUser(db *sqlx.DB, adminPassword *string) (int64, error) {
	var user model.User
	_ = db.Get(&user, `SELECT * FROM users WHERE username = ?`, AdminUsername)
	if user.ID > 0 {
		return user.ID, nil
	}

	if adminPassword == nil {
		defaultPwd := DefaultAdminPassword
		adminPassword = &defaultPwd
	}

	password, err := crypto.Encrypt(*adminPassword)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(`INSERT INTO users (username, email, role, name, password) 
		VALUES (?, ?, ?, ?, ?)`,
		AdminUsername, AdminEmail, AdminRole, AdminName, password)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
