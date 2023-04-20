package main

import (
	"encoding/json"
	"io/ioutil"
	"issue-api/middleware"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var issueMap = map[string]IssueGetModel{}
var issueList = []string{}

var apiPrefix string = "issue_api"
var create_issue string = apiPrefix + ":create_issue"
var list_issues string = apiPrefix + ":list_issues"
var close_issue string = apiPrefix + ":close_issue"

type IssueGetModel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type IssuePostModel struct {
	Name string `json:"name"`
}

type IssuePatchModel struct {
	Status string `json:"status"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/issues", GetIssues).Methods("GET", "OPTIONS")
	router.HandleFunc("/issues", CreateIssue).Methods("POST", "OPTIONS")
	router.HandleFunc("/issues/{id}", UpdateIssue).Methods("PATCH", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}

func GetIssues(w http.ResponseWriter, r *http.Request) {

	if middleware.Authorize(w, r, []string{list_issues}) {
		issueGetList := make([]IssueGetModel, 0, len(issueMap))
		for _, issue := range issueList {
			issueGetList = append(issueGetList, issueMap[issue])
		}

		json.NewEncoder(w).Encode(issueGetList)
	}
}

func CreateIssue(w http.ResponseWriter, r *http.Request) {

	if middleware.Authorize(w, r, []string{create_issue}) {

		var payload IssuePostModel
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &payload)

		id := uuid.New().String()
		log.Println(id)
		issueList = append(issueList, id)
		issueMap[id] = IssueGetModel{id, payload.Name, "Open"}

		json.NewEncoder(w).Encode(issueMap[id])
	}
}

func UpdateIssue(w http.ResponseWriter, r *http.Request) {

	if middleware.Authorize(w, r, []string{close_issue}) {

		var payload IssuePatchModel
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &payload)

		params := mux.Vars(r)
		id := params["id"]
		issue := issueMap[id]
		issueMap[id] = IssueGetModel{issue.ID, issue.Name, payload.Status}

		json.NewEncoder(w).Encode(issueMap[id])
	}
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			h.ServeHTTP(w, r)
		}
	}
}
