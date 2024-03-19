package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// 이것이 오류를 반환하면, error.As() 함수를 사용하여 오류가 *mysql.MySQLError 유형인지 확인합니다.
		// 그렇다면 오류는 mySQLError 변수에 할당됩니다. 그러면 우리는 할 수 있다
		// 오류 코드가 1062인지 확인하고 오류 내용이 users_uc_email 키와 관련되어 있는지 확인하세요.
		// 오류 메시지 문자열. 그렇다면 ErrDuplicateEmail 오류가 반환됩니다.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
