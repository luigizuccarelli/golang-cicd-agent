package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/schema"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/service"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func CICDHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var git *schema.GiteaSchema
	var mapping *schema.MapBinding

	body, err := ioutil.ReadAll(r.Body)
	con.Trace("Input data %s", string(body))
	if err != nil {
		con.Error("CICDHandler could not read body data %v", err)
		resp := "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\"" + fmt.Sprintf("CICDHandler could not read body data %v", err) + "\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, resp)
		return
	}

	err = json.Unmarshal(body, &git)
	if err != nil {
		con.Error("CICDHandler could not unmarshal to struct %v", err)
		resp := "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\"" + fmt.Sprintf("CICDHandler could not unmarshal struct %v", err) + "\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, resp)
		return
	}

	con.Trace("CICDHandler WEBHOOK_SECRET : %s : %s:", git.Secret, os.Getenv("WEBHOOK_SECRET"))
	apikey := strings.Trim(os.Getenv("WEBHOOK_SECRET"), "\n")
	// first check secret
	if git.Secret != apikey {
		con.Error("CICDHandler api secret invalid")
		resp := "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\"CICDHandler api secret invalid\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, resp)
		return
	}

	con.Debug("Mapping struct %v", git)

	// we now post to our various eventlisteners
	if git.Action == "published" || git.Action == "closed" {
		// only post on merged true
		if git.PullRequest.Merged {
			// post to the dev eventlistener - for our normal cicd dev build
			mapping = &schema.MapBinding{
				RepoUrl:    git.Repository.CloneURL,
				RepoName:   git.Repository.Name,
				RepoHash:   git.PullRequest.MergeCommitSha,
				ActorName:  git.PullRequest.User.Login,
				ActorEmail: git.PullRequest.User.Email,
				Message:    git.PullRequest.Title,
				Env:        "dev",
			}
		}

		if git.Action == "published" {
			// check for the prerelease field
			if git.Release.Prerelease {
				// post to the uat eventlistener
				mapping = &schema.MapBinding{
					RepoUrl:    git.Repository.CloneURL,
					RepoName:   git.Repository.Name,
					RepoHash:   git.Release.TargetCommitish,
					ActorName:  git.Release.Author.Login,
					ActorEmail: git.Release.Author.Email,
					Message:    git.Release.Name + " " + git.Release.Body,
					TagVersion: git.Release.TagName,
					Env:        "uat",
				}
			} else {
				// post to prod eventlistener
				mapping = &schema.MapBinding{
					RepoUrl:    git.Repository.CloneURL,
					RepoName:   git.Repository.Name,
					RepoHash:   git.Release.TargetCommitish,
					ActorName:  git.Release.Author.Login,
					ActorEmail: git.Release.Author.Email,
					Message:    git.Release.Name + " " + git.Release.Body,
					TagVersion: git.Release.TagName,
					Env:        "prod",
				}
			}
		}
		// ignore infra mapping for now
		/*
		infra, err := getInfraRepo(mapping.RepoName)
		if err != nil {
			resp := "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\"" + fmt.Sprintf(" %v", err) + "\"}"
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, resp)
			return
		}
		mapping.InfraRepo = infra
		*/
		// execute the CICD pipeline in a seperate light weight thread
		go service.ExecutePipeline(mapping, con)

		resp := "{\"status\":\"OK\", \"statuscode\":\"200\",\"message\":\"Pipeline started - view progress via cicd web console\"}"
		w.WriteHeader(http.StatusOK)
		con.Debug("Result struct for gitea webhook %v", mapping)
		fmt.Fprintf(w, string(resp))
	} else {
		con.Debug("NOP - no merge or release")
	}
}

// getInfraRepo - utility function to get the mapped infra repo
func getInfraRepo(name string) (string, error) {
	var result string

	repos := strings.Split(os.Getenv("REPO_MAPPING"), "\n")
	//prefix := strings.Split(name, "-")
	for x, _ := range repos {
		if strings.Contains(repos[x], name) {
			result = strings.Split(repos[x], "=")[1]
			break
		}
	}
	if result == "" {
		return "", errors.New("Infra repo not found")
	}
	return result, nil
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
	return
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("API-KEY") != "" {
		w.Header().Set("API_KEY_PT", r.Header.Get("API_KEY"))
	}
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
