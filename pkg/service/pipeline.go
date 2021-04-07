package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/schema"
)

// ExecutePipeline - says what it does :)
func ExecutePipeline(cicd *schema.MapBinding, con connectors.Clients) error {
	var pipeline []schema.Pipeline
	// always clone first
	cicd.Action = "git"
	_, err := taskHandler(cicd, con)
	if err != nil {
		return err
	}
	// check if the .cicd/pipeline.json file is present in the repo
	data, err := ioutil.ReadFile(os.Getenv("WORKDIR") + "/" + cicd.RepoName + "/.cicd/pipeline-" + cicd.Env + ".json")
	if err != nil {
		return err
	}
	// unmarshall the json file to struct
	err = json.Unmarshal(data, &pipeline)
	if err != nil {
		return err
	}

	// ensure tasks are ordered properly
	sort.Slice(pipeline[:], func(i, j int) bool {
		return pipeline[i].Id < pipeline[j].Id
	})

	for i := 0; i < len(pipeline); i++ {
		cicd.Action = pipeline[i].Task
		_, e := taskHandler(cicd, con)
		if e != nil {
			return e
		}
	}
	return nil
}

// scmMapper - map params for scm clone
func scmMapper(cicd *schema.MapBinding) []string {
	cicd.ActionDetail = "clone"
	cicd.Workdir = os.Getenv("WORKDIR")
	return []string{"clone", cicd.RepoUrl, cicd.RepoName}
}

// makeMapper - map params to execute make
func makeMapper(cicd *schema.MapBinding) []string {
	cicd.Workdir = os.Getenv("WORKDIR") + "/" + cicd.RepoName
	cicd.ActionDetail = cicd.Action
	cicd.Action = "make"
	return []string{cicd.ActionDetail}
}

// sonarMapper - map params to execute sonarqube code analysis
func sonarMapper(cicd *schema.MapBinding) []string {
	cicd.Workdir = os.Getenv("WORKDIR") + "/" + cicd.RepoName
	cicd.ActionDetail = "sonarscan"
	cicd.Action = os.Getenv("SONARSCAN_PATH") + "/sonar-scanner"
	args := []string{}
	return args
}

// sonarResultsMapper - map params to get quality gate info from sonarqube
func sonarResultsMapper(cicd *schema.MapBinding) []string {
	cicd.Workdir = os.Getenv("WORKDIR") + "/" + cicd.RepoName
	cicd.ActionDetail = "sonar-results"
	cicd.Action = "curl"
	args := []string{"-s", "-H", "Content-Type: application/json", "-H", "Accept: application/json", "-u", os.Getenv("SONAR_TOKEN") + ":", os.Getenv("SONAR_URL") + "/api/qualitygates/project_status?projectKey=" + cicd.RepoName}
	return args
}

// taskHandler - executes specific cicd tasks
func taskHandler(cicd *schema.MapBinding, con connectors.Clients) (string, error) {
	var args []string
	var err error
	var out string

	switch cicd.Action {
	case "git":
		_, e := os.Stat(os.Getenv("WORKDIR") + "/" + cicd.RepoName)
		// we always clean the directory before cloning
		if e == nil {
			out, err = con.ExecOS(os.Getenv("WORKDIR"), "rm", []string{"-rf", cicd.RepoName}, true)
		}
		// prepare console output for webconsole to use
		os.MkdirAll(os.Getenv("CICD_CONSOLE_DIR")+"/"+cicd.RepoName, os.ModePerm)
		if err != nil {
			con.Error("result : %s", out)
			con.Error("stderr : %v", err)
			return out, err
		}
		args = scmMapper(cicd)
		break
	case "clean":
		args = makeMapper(cicd)
		break
	case "cover":
		args = makeMapper(cicd)
		break
	case "test":
		args = makeMapper(cicd)
		break
	case "build":
		args = makeMapper(cicd)
		break
	case "sonarscan":
		props := "sonar.host.url=" + os.Getenv("SONAR_URL") + "\n" + "sonar.login=" + os.Getenv("SONAR_TOKEN") + "\n"
		f, _ := os.OpenFile(os.Getenv("WORKDIR")+"/"+cicd.RepoName+"/sonar-project.properties", os.O_APPEND|os.O_WRONLY, 0755)
		f.WriteString(props)
		f.Close()
		args = sonarMapper(cicd)
		break
	case "sonarresults":
		args = sonarResultsMapper(cicd)
		break
	case "sonaranalyse":
		data, _ := ioutil.ReadFile(os.Getenv("CICD_CONSOLE_DIR") + "/" + cicd.RepoName + "/sonaranalyse.txt")
		if strings.Contains(string(data), "{\"projectStatus\":{\"status\":\"OK\"") {
			args = []string{"-e", "PASS:"}
		} else {
			args = []string{"-e", "ERROR:"}
		}
		cicd.Action = "echo"
		break
	case "container":
		args = makeMapper(cicd)
		break
	case "push":
		args = makeMapper(cicd)
		break
	}
	out, err = con.ExecOS(cicd.Workdir, cicd.Action, args, true)
	if err != nil {
		// if the file write doesn't work we jus won't see any info on the webconsole - so we just ignore errors
		ioutil.WriteFile(os.Getenv("CICD_CONSOLE_DIR")+"/"+cicd.RepoName+"/"+cicd.ActionDetail+".txt", []byte("ERROR:</br>"+out+fmt.Sprintf("%v", err)), 0755)
	} else {
		if out == "" {
			out = "PASS"
		}
		// if the file write doesn't work we jus won't see any info on the webconsole - so we just ignore errors
		ioutil.WriteFile(os.Getenv("CICD_CONSOLE_DIR")+"/"+cicd.RepoName+"/"+cicd.ActionDetail+".txt", []byte("PASS:</br>"+out), 0755)
	}
	return out, err
}
