package main

import (
	"database/sql"
	"log"
)

type User struct {
	Uid      string `json:"uid"`
	Name     string `json:"name"`
	Picture  string `json:"pic"`
	Phone    string `json:"phone"`
	Otp      string `json:"otp"`
	Verified bool   `json:"verified"`
}

func saveUser(db *sql.DB, user User) error {
	var duplicate bool
	err := db.QueryRow("SELECT 1 from user WHERE uid = ?", user.Uid).Scan(&duplicate)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if !duplicate {
		_, err := db.Exec(`INSERT INTO user (uid, name, pic, phone, otp, verified) VALUES 
		(?, ?, ?, ?, ?, ?)`,
			user.Uid, user.Name, user.Picture, user.Phone, user.Otp, user.Verified)
		return err
	}

	log.Println("Non unique duplicate but still free")

	return nil
}

func updateUser(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE user SET name = ?, pic =? WHERE phone = ?",
		user.Name, user.Picture, user.Phone)

	return err
}

func getUser(db *sql.DB, phone string) (User, error) {
	sqlParam := phone
	row := db.QueryRow("SELECT * FROM user WHERE phone = ?", sqlParam)

	var user User
	err := row.Scan(&user.Uid, &user.Name, &user.Picture, &user.Phone, &user.Otp,
		&user.Verified)

	return user, err
}

func verifyUser(db *sql.DB, phone string, verified bool) error {
	_, err := db.Exec("UPDATE user SET verified = ? WHERE phone = ?", phone, verified)
	return err
}
