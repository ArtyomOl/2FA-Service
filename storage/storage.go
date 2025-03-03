package storage

import (
	"authentification_service/safety"
	"authentification_service/totp/random"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db         *sql.DB
	stmtInsert *sql.Stmt
	stmtSelect *sql.Stmt
}

// Подключение к базе данных
func New() (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open("mysql", "root:***@tcp(localhost:3306)/2fa_service")
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO users(login, password, code) VALUES(?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtSelect, err := db.Prepare("SELECT code FROM users WHERE login = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db, stmtInsert, stmtSelect}, nil
}

// Добавление данных нового пользователя в базу данных
func (s *Storage) SaveUser(login string, password string) error {
	const op = "storage.SaveUser"

	random_code := random.RandomToken()

	login_hash := safety.Hashing(login)
	password_hash := safety.Hashing(password)

	_, err := s.stmtInsert.Exec(login_hash, password_hash, random_code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println("successful safe")

	return err
}

// Получение текущего кода пользователя
func (s *Storage) GetCode(login string) (string, error) {
	const op = "storage.GetCode"

	var code string

	login_hash := safety.Hashing(login)

	err := s.stmtSelect.QueryRow(login_hash).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "-", nil
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return code, nil
}

// Закрытие подключения к бд
func (s *Storage) Close() error {
	const op = "storage.Close"
	errs := make([]string, 0, 3)
	if err := s.stmtInsert.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := s.stmtSelect.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := s.db.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) != 0 {
		return fmt.Errorf("%s: %s", op, strings.Join(errs, ", "))
	}
	return nil
}
