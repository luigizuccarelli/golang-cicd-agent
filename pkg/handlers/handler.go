package handlers

import (
	"encoding/json"
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
		fmt.Fprintf(w, "%s", resp)
		return
	}

	err = json.Unmarshal(body, &git)
	if err != nil {
		con.Error("CICDHandler could not unmarshal to struct %v", err)
		resp := "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\"" + fmt.Sprintf("CICDHandler could not unmarshal struct %v", err) + "\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", resp)
		return
	}

	con.Trace("CICDHandler WEBHOOK_SECRET : %s : %s:", git.Secret, os.Getenv("WEBHOOK_SECRET"))
	apikey := strings.Trim(os.Getenv("WEBHOOK_SECRET"), "\n")
	// first check secret
	if git.Secret != apikey {
		con.Error("CICDHandler api secret invalid")
		resp := "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\"CICDHandler api secret invalid\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", resp)
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

		// execute the CICD pipeline in a separate light weight thread
		go func() {
			err := service.ExecutePipeline(mapping, con)
			if err != nil {
				con.Error("Error executing pipeline %v", err)
			}
		}()

		resp := "{\"status\":\"OK\", \"statuscode\":\"200\",\"message\":\"Pipeline started - view progress via cicd web console\"}"
		w.WriteHeader(http.StatusOK)
		con.Debug("Result struct for gitea webhook %v", mapping)
		fmt.Fprintf(w, "%s", string(resp))
	} else {
		con.Debug("NOP - no merge or release")
	}
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
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
