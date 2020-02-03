package models

import (
	"database/sql"
	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	"log"
)

type User struct {
	Id           int
	Email        string
	PasswordHash string
	Admin        int
}

func (u *User) CreateNewUser() {
	var db = database.ConnectDB()
	defer db.Close()
	sqlStatement := `INSERT INTO users (Email, PasswordHash, Admin) VALUES ($1, $2, $3) RETURNING id`

	id := 0
	err := db.QueryRow(sqlStatement, u.Email, u.PasswordHash).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	u.Id = id
}

func (u *User) GetUserByEmail(email string) *User {
	var db = database.ConnectDB()
	defer db.Close()

	sqlStatement := `SELECT * FROM public.users WHERE Email=$1;`
	user := new(User)
	row := db.QueryRow(sqlStatement, email)
	err := row.Scan(user.Id, user.Email, user.PasswordHash, user.Admin)
	switch err {
	case sql.ErrNoRows:
		log.Printf("Unable to find user with email: %s", email)
		return nil
	case nil:
		return user
	default:
		panic(err)
	}
}
