#!/bin/bash

LOG_LEVEL=debug
SERVER_PORT=9001
VERSION=1.1.0
NAME=golang-cicd-webconsole
CICD_CONSOLE_DIR=../console
SONAR_TOKEN=194eebd4eb150a5041ac93e30be387c6729238d4
SONAR_URL=http://192.168.1.25:9000
SONARSCAN_PATH=/home/lzuccarelli/sonar-scanner-4.6.0.2311/bin
WORKDIR=../workdir
REPO_MAPPING=golang-gitwebhook-service=infra-golang-gitwebhook-service
WEBHOOK_SECRET=lmz123

export LOG_LEVEL NAME SERVER_PORT VERSION CICD_CONSOLE_DIR SONAR_TOKEN SONAR_URL SONARSCAN_PATH WORKDIR REPO_MAPPING WEBHOOK_SECRET

./build/microservice
