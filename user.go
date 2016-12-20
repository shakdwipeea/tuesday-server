package main

import "database/sql"

type User struct {
	Uid      int    `json:"uid"`
	Name     string `json:"name"`
	Picture  string `json:"pic"`
	Phone    string `json:"phone"`
	Otp      string `json:"otp"`
	Verified bool   `json:"verified"`
	Token    string `json:"token"`
}

func saveUser(db *sql.DB, user User) (int, error) {
	res, err := db.Exec(`INSERT INTO user (uid, name, pic, phone, otp, verified) VALUES 
		(?, ?, ?, ?, ?, ?)`,
		user.Uid, user.Name, user.Picture, user.Phone, user.Otp, user.Verified)

	lastID, err := res.LastInsertId()
	return int(lastID), err
}

func updateOtp(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE user SET otp = ?, phone = ? WHERE phone = ?",
		user.Otp, user.Verified, user.Phone)

	return err
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
