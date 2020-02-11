package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dimitarvalkanov7/chaoscamp-demo/models/demoscene"
	"github.com/dimitarvalkanov7/chaoscamp-demo/services/github"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/services/cmd"
)

var (
	imagePath = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/dockerimages/"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "demoscenes", "index.html")))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

func Repositories(w http.ResponseWriter, r *http.Request) {
	//res := docker.IsPortAvialable("1235")
	// rb := make([]RepoWithBranch, 0)
	// rb = append(rb, RepoWithBranch{Repo: "chaoscamp", Branch: "master"})
	// rb = append(rb, RepoWithBranch{Repo: "golangrulz", Branch: "master"})
	// CreateDemoscene("test-demo-scene", rb)
	res := github.GetBranchesByRepository()

	type Context struct {
		BranchesByRepo map[string][]string
	}
	data := Context{BranchesByRepo: res}

	//cmd.Execute("git clone -b v1 https://github.com/docker-training/node-bulletin-board /home/leron/DockerTutorial")
	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "demoscenes/repository.html")))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

type Demoscene struct {
	Name string
	Rb   []RepoWithBranch
}

type RepoWithBranch struct {
	Repo   string
	Branch string
}

type Exception struct {
	Data   string
	Status string
}

func CreateDemoscene(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("rb")
	ds := new(Demoscene)
	json.Unmarshal([]byte(data), &ds)

	if nameIsNotUnique(ds.Name) {
		msg := fmt.Sprintf("Name %s alredy exists.", ds.Name)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Exception{Data: msg, Status: "error"})
		return
	}

	buildDemoscene(ds.Name, ds.Rb)

	time.Sleep(2 * time.Second)
}

func nameIsNotUnique(sceneName string) bool {
	scene := demoscene.GetSceneByName(sceneName)
	if scene == nil {
		return false
	}

	return true
}

func buildDemoscene(sceneName string, repoBranch []RepoWithBranch) {
	p := path.Join(imagePath, sceneName)
	cloneCommand := "git clone -b {branch} --single-branch git@github.com:dimitarvalkanov7/{repo}.git " + p

	var wg sync.WaitGroup
	createNewDirectoryWithName(sceneName)
	var cc string
	for _, rb := range repoBranch {
		cc = strings.Replace(cloneCommand, "{branch}", rb.Branch, -1)
		cc = strings.Replace(cc, "{repo}", rb.Repo, -1)
		wg.Add(1)
		r := rb.Repo
		command := cc + "/" + r
		go func(command, r string) {
			defer wg.Done()
			createDir(p, r)
			cmd.Execute(command)
		}(command, r)
	}
	wg.Wait()

	var imageName string
	var pathToDockerfile string
	for _, rebr := range repoBranch {
		imageName = rebr.Repo + "_" + sceneName
		pathToDockerfile = path.Join(imagePath, sceneName, rebr.Repo)

		wg.Add(1)
		go func(iname, ptdf string) {
			defer wg.Done()
			pathWithDockerfile := path.Join(ptdf, "Dockerfile")
			com := fmt.Sprintf("docker build -t %s -f %s %s", iname, pathWithDockerfile, ptdf)
			cmd.Execute(com)
		}(imageName, pathToDockerfile)
	}
	wg.Wait()

	//save to database
	for _, rb := range repoBranch {
		demoscene.CreateNewScene(sceneName, rb.Repo, rb.Branch)
	}
}

func createNewDirectoryWithName(dirName string) {
	os.Mkdir(path.Join(imagePath, dirName), os.FileMode(0777))
}

func createDir(pathToDir, dirName string) {
	os.Mkdir(path.Join(pathToDir, dirName), os.FileMode(0777))
}
