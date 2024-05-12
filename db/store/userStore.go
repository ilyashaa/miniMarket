package store

import (
	"database/sql"
	"fmt"
	"log"
	"miniMarket/product"
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

func (store *UserStore) Close() {
	if store.db != nil {
		store.db.Close()
	}
}

func (store *UserStore) RegisterSQL(email string, passwordHash string) (string, error) {

	db := store.db

	query := `
    INSERT INTO users (email, password)
    VALUES ($1, $2);`

	result, err := db.Exec(query, email, passwordHash)
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
	if sqlBirthdayDate == "" {
		sqlBirthdayDate = "01.06.2000"
	}

	return sqlEmail, sqlNickname, sqlBirthdayDate, nil
}

func (store *UserStore) UpdateInfoSQL(email, nickname, birthdayDate string) {
	db := store.db

	query := `
    UPDATE users SET nickname = $1, birthdayDate = $2 WHERE email = $3`

	result, err := db.Exec(query, nickname, birthdayDate, email)
	if err != nil {
		// return "Не получилось передать данные на сервер!", err
	}

	_, err = result.RowsAffected()
	if err != nil {
		// return "Не получилось передать данные на сервер.", err
	}
}

func (store *UserStore) GetProduct() ([]product.Product, error) {

	db := store.db
	rows, err := db.Query("SELECT name, price, image FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []product.Product
	for rows.Next() {
		var p product.Product
		if err := rows.Scan(&p.Name, &p.Price, &p.Image); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}
