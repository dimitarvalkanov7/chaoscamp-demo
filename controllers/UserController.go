package controllers

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	"github.com/dimitarvalkanov7/chaoscamp-demo/encryption"
	"github.com/dimitarvalkanov7/chaoscamp-demo/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/"
)

var db = database.ConnectDB()

func Login(w http.ResponseWriter, r *http.Request) {
	_, err := encryption.GetLoggedUser(r)
	if err == nil {
		log.Println("User not found.")
		http.Redirect(w, r, "/auth/home", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates/user", "login.html")))
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}

}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	loggedUser := FindOne(email, password)

	if loggedUser == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	encEmail := encryption.Encrypt(loggedUser.Email)

	cookie := http.Cookie{Name: "demoscenes", Value: encEmail, Expires: time.Now().Add(24 * time.Hour), Path: "/"}
	r.AddCookie(&cookie)
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/auth/home", 302)
}

func FindOne(email, password string) *models.User {
	var user *(models.User)
	user = user.GetUserByEmail(email)
	if user == nil {
		return nil
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil
	}

	return user
}

func SeedUsers() error {
	if FindOne("admin@gmail.com", "password") != nil {
		return nil
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := new(models.User)
	u.Email = "admin@gmail.com"
	u.PasswordHash = string(pwd)
	u.IsAdmin = 1
	u.IsVerified = 1
	u.CreateNewUser()

	return nil
}
