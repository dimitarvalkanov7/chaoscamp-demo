package controllers

import (
	//"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	"github.com/dimitarvalkanov7/chaoscamp-demo/encryption"
	"github.com/dimitarvalkanov7/chaoscamp-demo/models"
	"github.com/dimitarvalkanov7/chaoscamp-demo/services/github"
	"golang.org/x/crypto/bcrypt"
)

const (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/"
)

var db = database.ConnectDB()

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := encryption.GetLoggedUser(r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}

	// TODO For test only
	//github.GetBranches()
	res := github.GetBranchesByRepository()
	//res := github.GetAllRepositories()
	type Context struct {
		BranchesByRepo map[string][]string
	}
	data := Context{BranchesByRepo: res}
	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "demoscenes", "index.html")))
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

func AddNewUser(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates/user", "new.html")))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	// c, err := r.Cookie("demoscenes")
	// if err == nil {
	// 	if c.Expires.Before(time.Now()) && len(c.Value) > 3 {
	// 		http.Redirect(w, r, "/demo-scenes", 302)
	// 		return
	// 	}
	// }

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
		u.IsAdmin = 0
		u.IsVerified = 0
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
	_, err := encryption.GetLoggedUser(r)
	if err == nil {
		log.Println(err)
		http.Redirect(w, r, "/home", 302)
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
	// TODO decide what to do on failure
	if loggedUser == nil {
		log.Fatal("Unable to login")
		http.Redirect(w, r, "/login", 302)
		return
	}

	// TODO decide what to do on success
	encEmail := encryption.Encrypt(loggedUser.Email)

	cookie := http.Cookie{Name: "demoscenes", Value: encEmail, Expires: time.Now().Add(24 * time.Hour)}
	r.AddCookie(&cookie)
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/home", 302)
}

func FindOne(email, password string) *models.User {
	var user *(models.User)
	user = user.GetUserByEmail(email)
	if user == nil {
		return nil
	}

	//expiresAt := time.Now().Add(time.Minute * 100000).Unix()

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
