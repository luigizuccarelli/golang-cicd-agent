#!/bin/bash

LOG_LEVEL=trace
SERVER_PORT=9004
VERSION=1.1.0
NAME=golang-cicd-agent
CICD_CONSOLE_DIR=/home/lzuccarelli/cicd/console
SONAR_TOKEN=365a52151abce2ebfe5855fd591996ab3dedb133
SONAR_URL=http://192.168.8.118:9000
SONARSCAN_PATH=/home/lzuccarelli/Programs/sonar-scanner-4.6.0.2311-linux/bin
WORKDIR=/home/lzuccarelli/cicd/workdir
REPO_MAPPING=golang-gitwebhook-service=infra-golang-gitwebhook-service
WEBHOOK_SECRET=lmz123

export LOG_LEVEL NAME SERVER_PORT VERSION CICD_CONSOLE_DIR SONAR_TOKEN SONAR_URL SONARSCAN_PATH WORKDIR REPO_MAPPING WEBHOOK_SECRET

./build/agent
