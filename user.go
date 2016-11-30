package main

import (
	"database/sql"
	"log"
)

type User struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Picture string `json:"pic"`
}

func saveUser(db *sql.DB, user User) error {
	var duplicate bool
	err := db.QueryRow("SELECT 1 from user WHERE uid = ?", user.Uid).Scan(&duplicate)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if !duplicate {
		_, err := db.Exec("INSERT INTO user (uid, name, pic) VALUES (?, ?, ?)",
			user.Uid, user.Name, user.Picture)
		return err
	}

	log.Println("Non unique duplicate but still free")

	return nil
}

func getUsers(db *sql.DB, prefix string) ([]User, error) {
	var users []User

	sqlParam := "%" + prefix + "%"
	rows, err := db.Query("SELECT * FROM user WHERE name like ?", sqlParam)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Uid, &user.Name, &user.Picture); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}
