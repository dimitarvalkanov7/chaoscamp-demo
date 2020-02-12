package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	USER        = "dimitarvalkanov7"
	API_TOKEN   = "6f731f27e7aa2962e3959d0d97b63944c0f1a7e1"
	GIT_API_URL = "https://api.github.com/users/dimitarvalkanov7/repos"
	fullPath    = "https://api.github.com/users/dimitarvalkanov7/repos?access_token=04d2e2ff594e6dba68ffefe03e0af6871351e8af"
)

var (
	usedRepositories = []string{"chaoscamp", "golangrulz"}
)

func GetAllRepositories() []Repository {
	client := &http.Client{}
	token := fmt.Sprintf("%s/token:%s", USER, API_TOKEN)

	currTokenLen := len(token)
	for {
		token = strings.Replace(token, "\n", "", -1)
		if currTokenLen == len(token) {
			break
		}
		currTokenLen = len(token)
	}
	encodedToken := encodeBase64([]byte(token))

	request, err := http.NewRequest("GET", GIT_API_URL, nil)
	authHeader := fmt.Sprintf("Basic %s", encodedToken)
	if err != nil {
		log.Println(err)
		return nil
	}

	// add header to request
	request.Header.Set("Authorization", authHeader)

	// perform get request
	var text string
	res, err := client.Do(request)
	if err != nil {
		log.Println(err)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		text = string(body)
	}

	repositories := make([]Repository, 0)
	json.Unmarshal([]byte(text), &repositories)

	// repoNames := make([]string, len(repositories))
	// for _, v := range repositories {
	// 	repoNames = append(repoNames, v.Name)
	// }

	// return repoNames
	return repositories
}

func GetAllBranches() {
	client := &http.Client{}
	token := fmt.Sprintf("%s/token:%s", USER, API_TOKEN)

	currTokenLen := len(token)
	for {
		token = strings.Replace(token, "\n", "", -1)
		if currTokenLen == len(token) {
			break
		}
		currTokenLen = len(token)
	}
	encodedToken := encodeBase64([]byte(token))

	//request, err := http.NewRequest("GET", GIT_API_URL, nil)
	branches_url := "https://api.github.com/repos/dimitarvalkanov7/chaoscamp-demo/branches"
	request, err := http.NewRequest("GET", branches_url, nil)
	authHeader := fmt.Sprintf("Basic %s", encodedToken)
	if err != nil {
		log.Println(err)
		return
	}

	// add header to request
	request.Header.Set("Authorization", authHeader)

	// perform get request
	var text string
	res, err := client.Do(request)
	if err != nil {
		log.Println(err)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		text = string(body)
		log.Println(text)
	}
}

func GetAllBranchesForRepository(repository Repository) []Branch {
	client := &http.Client{}
	token := fmt.Sprintf("%s/token:%s", USER, API_TOKEN)

	currTokenLen := len(token)
	for {
		token = strings.Replace(token, "\n", "", -1)
		if currTokenLen == len(token) {
			break
		}
		currTokenLen = len(token)
	}
	encodedToken := encodeBase64([]byte(token))

	//request, err := http.NewRequest("GET", GIT_API_URL, nil)
	//branches_url := "https://api.github.com/repos/dimitarvalkanov7/chaoscamp-demo/branches"
	branches_url := fmt.Sprintf("https://api.github.com/repos/dimitarvalkanov7/%s/branches", repository.Name)
	request, err := http.NewRequest("GET", branches_url, nil)
	authHeader := fmt.Sprintf("Basic %s", encodedToken)
	if err != nil {
		log.Println(err)
		return nil
	}

	// add header to request
	request.Header.Set("Authorization", authHeader)

	// perform get request
	var text string
	res, err := client.Do(request)
	if err != nil {
		log.Println(err)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		text = string(body)
	}

	branches := make([]Branch, 0)
	json.Unmarshal([]byte(text), &branches)

	return branches
}

func GetBranchesByRepository() map[string][]string {
	repositories := GetAllRepositories()

	branchesByRepository := make(map[Repository][]Branch)

	for _, repo := range repositories {
		for _, repoWithDockerFile := range usedRepositories {

			if repo.Name == repoWithDockerFile {
				branchesByRepository[repo] = GetAllBranchesForRepository(repo)
			}
		}
	}

	normalized := normalizeBranchRepoMap(branchesByRepository)

	return normalized
}

func normalizeBranchRepoMap(bbr map[Repository][]Branch) map[string][]string {
	normalized := make(map[string][]string)

	for repo, branches := range bbr {
		for _, branch := range branches {
			normalized[repo.Name] = append(normalized[repo.Name], branch.Name)
		}
	}

	return normalized
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

type Repository struct {
	Id                int64
	Node_id           string
	Name              string
	Full_name         string
	Private           bool
	Owner             Owner
	Html_url          string
	Description       string
	Fork              bool
	Url               string
	Forks_url         string
	Keys_url          string
	Collaborators_url string
	Teams_url         string
	Hooks_url         string
	Issue_events_url  string
	Events_url        string
	Assignees_url     string
	Branches_url      string
	Tags_url          string
	Blobs_url         string
	Git_tags_url      string
	Git_refs_url      string
	Trees_url         string
	Statuses_url      string
	Languages_url     string
	Stargazers_url    string
	Contributors_url  string
	Subscribers_url   string
	Subscription_url  string
	Commits_url       string
	Git_commits_url   string
	Comments_url      string
	Issue_comment_url string
	Contents_url      string
	Compare_url       string
	Merges_url        string
	Archive_url       string
	Downloads_url     string
	Issues_url        string
	Pulls_url         string
	Milestones_url    string
	Notifications_url string
	Labels_url        string
	Releases_url      string
	Deployments_url   string
	Created_at        string
	Updated_at        string
	Pushed_at         string
	Git_url           string
	Ssh_url           string
	Clone_url         string
	Svn_url           string
	Homepage          string
	Size              int
	Stargazers_count  int
	Watchers_count    int
	Language          string
	Has_issues        bool
	Has_projects      bool
	Has_downloads     bool
	Has_wiki          bool
	Has_pages         bool
	Forks_count       int
	Mirror_url        string
	Archived          bool
	Disabled          bool
	Open_issues_count int
	License           string
	Forks             int
	Open_issues       int
	Watchers          int
	Default_branch    string
	Permissions       Permissions
}

type Owner struct {
	Login               string
	Id                  int
	Node_id             string
	Avatar_url          string
	Gravatar_id         string
	Url                 string
	Html_url            string
	Followers_url       string
	Following_url       string
	Gists_url           string
	Starred_url         string
	Subscriptions_url   string
	Organizations_url   string
	Repos_url           string
	Events_url          string
	Received_events_url string
	Type                string
	Site_admin          bool
}

type Permissions struct {
	Admin bool
	Push  bool
	Pull  bool
}

type Branch struct {
	Name           string
	Commit         Commit
	Protected      bool
	Protection     Protection
	Protection_url string
}

type Commit struct {
	Sha string
	Url string
}

type Protection struct {
	Enabled                bool
	Required_status_checks Required_status_checks
}

type Required_status_checks struct {
	Enforcement_level string
	contexts          []string
}
