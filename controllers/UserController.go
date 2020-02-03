package controllers

import (
	"encoding/json"
	"github.com/dimitarvalkanov7/chaoscamp-demo/encryption"
	"html/template"
	"log"
	"net/http"
	"path"
	//"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	"github.com/dimitarvalkanov7/chaoscamp-demo/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/"
)

var db = database.ConnectDB()

func Register(w http.ResponseWriter, r *http.Request) {
	// c, err := r.Cookie("demoscenes")
	// if err == nil {
	// 	if c.Expires.Before(time.Now()) && len(c.Value) > 3 {
	// 		http.Redirect(w, r, "/demo-scenes", 302)
	// 		return
	// 	}
	// }

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "register.html")))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing template: %v\n", err)
		}
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	pass := r.FormValue("password")

	if email != "" && pass != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
		}

		u := new(models.User)
		u.Email = email
		u.PasswordHash = string(pwd)
		u.Admin = 0
		u.CreateNewUser()

		// http.SetCookie(w, &http.Cookie{
		// 	Name:  "demoscenes",
		// 	Value: string(u.Id),
		// })

		// http.Redirect(w, r, "/demo-scenes", 302)
		// return
	}
	// http.Redirect(w, r, "/register", 302)
	// return
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	loggedUser := FindOne(email, password)
	// TODO decide what to do on failure
	if loggedUser == nil {
		log.Fatal("Unable to login")
	}

	// TODO decide what to do on success
	userAsJson, _ := json.Marshal(loggedUser)
	encUser := encryption.Encrypt(string(userAsJson))

	http.SetCookie(w, &http.Cookie{
		Name:  "demoscenes",
		Value: encUser,
	})
}

func FindOne(email, password string) *models.User {
	var user *(models.User)
	user = user.GetUserByEmail(email)

	//expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil
	}

	return user
}
