# A simple golang microservice for FAST-CICD FE

Uses a simple post endpoint to render custom cicd results.

## Usage

start agent, webserver, webpipeline services

```bash
cd golang-cicd-agent
./run.sh &

cd ../golang-cicd-webpipeline
./run.sh &

cd ../golang-cicd-webserver
./run.sh &

```

to execute a pipeline

```bash
cd ../golang-cicd-agent
sed -i s/golang-repo-name/your-golang-repo-name-here/g tests/use-to-simulate.json 
# i.e sed -i s/golang-repo-name/golang-cicd-agent/g tests/use-to-simulate.json
curl -v -d'@tests/use-to-simulate.json' http://127.0.0.1:9004/api/v1/webhook
```

## Dashboard

View pipeline runs and results at 

```bash
http://127.0.0.1:9003/cicd-golang-cicd-agent-template.html
```
