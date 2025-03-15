package storage

import (
	"authentification_service/safety"
	"authentification_service/totp/random"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Storage struct {
	db         *sql.DB
	stmtInsert *sql.Stmt
	stmtSelect *sql.Stmt
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "***"
	dbname   = "2fa_service"
)

// Подключение к базе данных
func New() (*Storage, error) {
	const op = "storage.New"

	if password == "***" {
		fmt.Println("Change password")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtInsert, err := db.Prepare("INSERT INTO users(user_login, user_password, code) VALUES($1, $2, $3)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtSelect, err := db.Prepare("SELECT code FROM users WHERE user_login = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db, stmtInsert, stmtSelect}, nil
}

// Добавление данных нового пользователя в базу данных
func (s *Storage) SaveUser(login string, password string) error {
	fmt.Println("work!!!!!!!!!!!!!!!!!!!1")
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
func (s *Storage) GetCode(login string) (bool, string, error) {
	const op = "storage.GetCode"

	var code string

	login_hash := safety.Hashing(login)

	err := s.stmtSelect.QueryRow(login_hash).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", fmt.Errorf("%s: %w", op, err)
	}

	return true, code, nil
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
