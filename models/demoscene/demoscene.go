package demoscene

import (
	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	"log"
)

type Demoscene struct {
	Id             int
	Name           string
	RepositoryName string
	BranchName     string
	Port           int
	IsDeleted      int
}

func CreateNewScene(sceneName, repoName, branchName string, port int) Demoscene {
	scene := Demoscene{
		Name:           sceneName,
		RepositoryName: repoName,
		BranchName:     branchName,
		Port:           port,
		IsDeleted:      0,
	}
	var db = database.ConnectDB()
	defer db.Close()
	sqlStatement := `INSERT INTO demoscenes (Name, RepositoryName, BranchName, Port, IsDeleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	id := 0
	err := db.QueryRow(sqlStatement, sceneName, repoName, branchName, port, 0).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	scene.Id = id
	return scene
}

func GetSavedPorts() []int {
	var db = database.ConnectDB()
	defer db.Close()
	sqlStatement := "select port from demoscenes where isdeleted = 0"

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	ports := make([]int, 0)
	for rows.Next() {
		var p int
		err := rows.Scan((&p))
		if err != nil {
			log.Println(err)
			return nil
		}
		ports = append(ports, p)
	}
	return ports
}

func GetAll() []Demoscene {
	var db = database.ConnectDB()
	defer db.Close()

	sqlStatement := `SELECT Id, Name, RepositoryName, BranchName, Port, IsDeleted FROM public.demoscenes WHERE IsDeleted=0;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	scenes := make([]Demoscene, 0)
	for rows.Next() {
		scene := new(Demoscene)
		err := rows.Scan((&scene.Id), (&scene.Name), (&scene.RepositoryName), (&scene.BranchName), (&scene.Port), (&scene.IsDeleted))
		if err != nil {
			log.Println(err)
			return nil
		}
		scenes = append(scenes, *scene)
	}

	return scenes
}

func GetSceneByName(sceneName string) []Demoscene {
	var db = database.ConnectDB()
	defer db.Close()

	sqlStatement := `SELECT Id, Name, RepositoryName, BranchName, Port, IsDeleted FROM public.demoscenes WHERE Name=$1;`
	rows, err := db.Query(sqlStatement, sceneName)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	scenes := make([]Demoscene, 0)
	for rows.Next() {
		scene := new(Demoscene)
		err := rows.Scan((&scene.Id), (&scene.Name), (&scene.RepositoryName), (&scene.BranchName), (&scene.Port), (&scene.IsDeleted))
		if err != nil {
			log.Println(err)
			return nil
		}
		scenes = append(scenes, *scene)
	}

	return scenes
}

func DeleteDemosceneByName(sceneName string) {
	var db = database.ConnectDB()
	defer db.Close()
	sqlStatement := `UPDATE public.demoscenes SET Isdeleted=1 WHERE name=$1`

	_, err := db.Exec(sqlStatement, sceneName)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
