package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dimitarvalkanov7/chaoscamp-demo/models/demoscene"
	"github.com/dimitarvalkanov7/chaoscamp-demo/services/docker"
	"github.com/dimitarvalkanov7/chaoscamp-demo/services/github"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/services/cmd"
)

var (
	imagePath = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/dockerimages/"
)

func Home(w http.ResponseWriter, r *http.Request) {
	scene := make([]demoscene.Demoscene, 0)
	scene = demoscene.GetAll()
	s := make([]string, 0)
	for k, v := range scene {
		if k == 0 {
			s = append(s, v.Name)
		} else {
			contains := false
			for _, sc := range s {
				if sc == v.Name {
					contains = true
				}
			}
			if !contains {
				s = append(s, v.Name)
			}
		}
	}

	type SceneNames struct {
		Scenes []string
	}
	data := SceneNames{Scenes: s}

	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "demoscenes", "index.html")))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

func Repositories(w http.ResponseWriter, r *http.Request) {
	res := github.GetBranchesByRepository()

	type Context struct {
		BranchesByRepo map[string][]string
	}
	data := Context{BranchesByRepo: res}

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
	if len(scene) == 0 {
		return false
	}

	return true
}

type SceneContext struct {
	SceneName string
	Scene     []demoscene.Demoscene
}

func GetDemoscene(w http.ResponseWriter, r *http.Request) {
	sceneName := r.URL.Query().Get("name")
	scene := make([]demoscene.Demoscene, 0)
	scene = demoscene.GetSceneByName(sceneName)

	if scene == nil {
		log.Printf("Unable to find scene: %s", sceneName)
		http.Redirect(w, r, "/home", 302)
		return
	}

	data := SceneContext{SceneName: scene[0].Name, Scene: scene}

	runContainersForDemoscene(scene)

	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "demoscenes/demoscene.html")))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
}

func runContainersForDemoscene(scene []demoscene.Demoscene) {
	var wg sync.WaitGroup
	var currName string
	for _, s := range scene {
		wg.Add(1)
		currName = s.RepositoryName + "_" + s.Name
		go func(cName string) {
			defer wg.Done()
			stopContainer(cName)
		}(currName)
	}
	wg.Wait()

	command := "docker container run --name {name} -d -p {port}:8080 {originalName}"
	for _, s1 := range scene {
		currName = s1.RepositoryName + "_" + s1.Name
		c := strings.Replace(command, "{name}", currName, -1)
		c = strings.Replace(c, "{originalName}", currName, -1)
		c = strings.Replace(c, "{port}", strconv.Itoa(s1.Port), -1)

		wg.Add(1)
		go func(com string) {
			defer wg.Done()
			cmd.Execute(com)
		}(c)
	}
	wg.Wait()
}

func DeleteDemoscene(w http.ResponseWriter, r *http.Request) {
	sceneName := r.URL.Query().Get("name")
	containers := make([]string, 0)

	scene := demoscene.GetSceneByName(sceneName)
	for _, v := range scene {
		sn := v.RepositoryName + "_" + sceneName
		containers = append(containers, sn)
	}

	// stop containers
	var wg sync.WaitGroup
	for _, container := range containers {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			c1 := "docker stop " + c
			cmd.Execute(c1)
		}(container)
	}
	wg.Wait()

	// remove containers
	for _, container := range containers {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			c1 := "docker rm " + c
			cmd.Execute(c1)
		}(container)
	}
	wg.Wait()

	// remove images
	command := "docker image rm {image}"
	for _, image := range containers {
		wg.Add(1)
		c := strings.Replace(command, "{image}", image, -1)
		go func(com string) {
			defer wg.Done()
			cmd.Execute(com)
		}(c)
	}
	wg.Wait()

	// soft delete the scene from the DB
	demoscene.DeleteDemosceneByName(sceneName)

	// remove directory
	deleteSceneDirectory(sceneName)
}

func deleteSceneDirectory(sceneName string) {
	fullPath := path.Join(imagePath, sceneName)
	command := "rm -rf {path}"
	c := strings.Replace(command, "{path}", fullPath, -1)
	cmd.Execute(c)
}

func stopContainer(containerName string) {
	c1 := "docker stop " + containerName
	cmd.Execute(c1)
	c2 := "docker container rm " + containerName
	cmd.Execute(c2)
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
	var com string
	for _, rebr := range repoBranch {
		imageName = rebr.Repo + "_" + sceneName
		pathToDockerfile = path.Join(imagePath, sceneName, rebr.Repo)
		pathWithDockerfile := path.Join(pathToDockerfile, "Dockerfile")
		com = fmt.Sprintf("docker build -t %s -f %s %s", imageName, pathWithDockerfile, pathToDockerfile)

		wg.Add(1)
		go func(com string) {
			defer wg.Done()
			cmd.Execute(com)
		}(com)
	}
	wg.Wait()

	//save to database
	savedPorts := make([]int, 0)
	savedPorts = demoscene.GetSavedPorts()
	for _, rb := range repoBranch {
		var portToUse int
		for i := 50000; i < 60000; i++ {
			isUsed := false
			for _, v := range savedPorts {
				if v == i {
					isUsed = true
					break
				}
			}
			if !isUsed && docker.IsPortAvialable(strconv.Itoa(i)) {
				portToUse = i
				savedPorts = append(savedPorts, i)
				break
			}
		}
		demoscene.CreateNewScene(sceneName, rb.Repo, rb.Branch, portToUse)
	}
}

func createNewDirectoryWithName(dirName string) {
	os.Mkdir(path.Join(imagePath, dirName), os.FileMode(0777))
}

func createDir(pathToDir, dirName string) {
	os.Mkdir(path.Join(pathToDir, dirName), os.FileMode(0777))
}
