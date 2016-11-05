package main

import (
	"database/sql"
)

type User struct {
	Uid string `json:"uid"`
	Name string `json:"name"`
	Picture string `json:"pic"`
}

func saveUser(db *sql.DB, user User) error {
	_, err := db.Exec("INSERT INTO user (uid, name, pic) VALUES (?, ?, ?)",
		user.Uid, user.Name, user.Picture)
	return err
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
