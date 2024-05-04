package userDB

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func StartSQL() {
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

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterSQL(email string, passwordHash string) (string, error) {

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

func AuthorizeSQL(email string, password string) (string, error) {
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

func SelectInfoSQL(emailKey string) (string, string, string, error) {
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
