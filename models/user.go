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
	IsAdmin      int
	IsVerified   int
}

func (u *User) CreateNewUser() {
	var db = database.ConnectDB()
	defer db.Close()
	sqlStatement := `INSERT INTO users (Email, PasswordHash, IsAdmin, IsVerified) VALUES ($1, $2, $3, $4) RETURNING id`

	id := 0
	err := db.QueryRow(sqlStatement, u.Email, u.PasswordHash, u.IsAdmin, u.IsVerified).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	u.Id = id
}

func (u *User) GetUserByEmail(email string) *User {
	var db = database.ConnectDB()
	defer db.Close()

	sqlStatement := `SELECT Id, Email, PasswordHash, IsAdmin, IsVerified FROM public.users WHERE Email=$1;`
	user := new(User)
	row := db.QueryRow(sqlStatement, email)
	err := row.Scan((&user.Id), (&user.Email), (&user.PasswordHash), (&user.IsAdmin), (&user.IsVerified))
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
