package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/dimitarvalkanov7/chaoscamp-demo/services/cmd"
)

const (
	imagePath = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/dockerimages/"
)

func Repositories(w http.ResponseWriter, r *http.Request) {
	//res := docker.IsPortAvialable("1235")
	rb := make([]RepoWithBranch, 0)
	rb = append(rb, RepoWithBranch{Repo: "chaoscamp", Branch: "master"})
	rb = append(rb, RepoWithBranch{Repo: "golangrulz", Branch: "master"})
	CreateDemoscene("test-demo-scene", rb)

	//cmd.Execute("git clone -b v1 https://github.com/docker-training/node-bulletin-board /home/leron/DockerTutorial")
	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "repository.html")))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

type RepoWithBranch struct {
	Repo   string
	Branch string
}

func CreateDemoscene(sceneName string, repoBranch []RepoWithBranch) {
	p := path.Join(imagePath, "test-demo-scene")
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
		//buildImage()
		//runContainer()
	}
	wg.Wait()
}

func createNewDirectoryWithName(dirName string) {
	os.Mkdir(path.Join(imagePath, dirName), os.FileMode(0777))
}

func createDir(pathToDir, dirName string) {
	os.Mkdir(path.Join(pathToDir, dirName), os.FileMode(0777))
}
