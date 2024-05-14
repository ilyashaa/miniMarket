package store

import (
	"database/sql"
)

func RegisterSQL(email string, passwordHash string, db *sql.DB) (string, error) {

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

func AuthorizeSQL(email string, password string, db *sql.DB) (string, error) {

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

func SelectInfoSQL(emailKey string, db *sql.DB) (string, string, string, error) {

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

func UpdateInfoSQL(email, nickname, birthdayDate string, db *sql.DB) {

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

// func GetProduct(db *sql.DB) ([]product.Product, error) {

// 	rows, err := db.Query("SELECT name, price, image FROM products")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var products []product.Product
// 	for rows.Next() {
// 		var p product.Product
// 		if err := rows.Scan(&p.Name, &p.Price, &p.Image); err != nil {
// 			return nil, err
// 		}
// 		products = append(products, p)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return products, nil

// }
