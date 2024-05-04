package store

import (
	"database/sql"
	"fmt"
	"log"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore() *UserStore {
	const (
		host     = "localhost"
		port     = 5432
		user     = "admin"
		password = "12345"
		dbname   = "postgres"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &UserStore{db: db}
}

func (store *UserStore) RegisterSQL(email string, passwordHash string) (string, error) {

	db := store.db

	defer db.Close()

	sqlStatement := `
    INSERT INTO users (email, password)
    VALUES ($1, $2);`

	result, err := db.Exec(sqlStatement, email, passwordHash)
	if err != nil {
		return "Не получилось передать данные на сервер!", err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return "Не получилось передать данные на сервер.", err
	}

	return "Вы прошли регистрацию, " + email, nil
}

func (store *UserStore) AuthorizeSQL(email string, password string) (string, error) {

	db := store.db

	defer db.Close()

	query := `SELECT Email, Password FROM users WHERE email = $1`

	rows, err := db.Query(query, email)
	if err != nil {
		return "Не удалось получить данные", err
	}

	defer rows.Close()
	var sqlEmail string
	var sqlPassword string

	for rows.Next() {
		err = rows.Scan(&sqlEmail, &sqlPassword)
		if err != nil {
			return "Не удалось расшифровать данные", err
		}

	}
	return sqlPassword, nil
}

func (store *UserStore) SelectInfoSQL(emailKey string) (string, string, string, error) {

	db := store.db

	defer db.Close()

	query := `SELECT email, nickname, birthdaydate FROM users WHERE email = $1`

	rows, err := db.Query(query, emailKey)
	if err != nil {
		return "", "", "", err
	}

	defer rows.Close()
	var sqlEmail, sqlNickname, sqlBirthdayDate string

	for rows.Next() {
		err = rows.Scan(&sqlEmail, &sqlNickname, &sqlBirthdayDate)
		if err != nil {
			return "", "", "", err
		}
	}

	return sqlEmail, sqlNickname, sqlBirthdayDate, nil
}
