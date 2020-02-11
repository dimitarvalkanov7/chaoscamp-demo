package demoscene

import (
	"database/sql"
	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	"log"
)

type Demoscene struct {
	Id             int
	Name           string
	RepositoryName string
	BranchName     string
	IsDeleted      int
}

func CreateNewScene(sceneName, repoName, branchName string) Demoscene {
	scene := Demoscene{
		Name:           sceneName,
		RepositoryName: repoName,
		BranchName:     branchName,
		IsDeleted:      0,
	}
	var db = database.ConnectDB()
	defer db.Close()
	sqlStatement := `INSERT INTO demoscenes (Name, RepositoryName, BranchName, IsDeleted) VALUES ($1, $2, $3, $4) RETURNING id`

	id := 0
	err := db.QueryRow(sqlStatement, sceneName, repoName, branchName, 0).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	scene.Id = id
	return scene
}

func GetSceneByName(sceneName string) *Demoscene {
	var db = database.ConnectDB()
	defer db.Close()

	sqlStatement := `SELECT Id, Name, RepositoryName, BranchName, IsDeleted FROM public.demoscenes WHERE Name=$1;`
	scene := new(Demoscene)
	row := db.QueryRow(sqlStatement, sceneName)
	err := row.Scan((&scene.Id), (&scene.Name), (&scene.RepositoryName), (&scene.BranchName), (&scene.IsDeleted))
	switch err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return scene
	default:
		panic(err)
	}
}
