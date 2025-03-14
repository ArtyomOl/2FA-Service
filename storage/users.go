package storage

import (
	"authentification_service/safety"
	"database/sql"
	"fmt"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func (s *Storage) CheckUser(login string, password string) (bool, error) {
	const op = "storage.users.CheckUser"

	stmtSelectAll, err := s.db.Prepare("SELECT user_id FROM users WHERE user_login = $1 AND user_password = $2")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var user_id int

	login_hash := safety.Hashing(login)
	password_hash := safety.Hashing(password)

	err = stmtSelectAll.QueryRow(login_hash, password_hash).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
