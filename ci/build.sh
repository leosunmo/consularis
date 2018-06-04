#!/bin/bash -ex

APP_NAME="consularis"
DOCKER_REPO=leosunmo/consularis
BRANCH=${TRAVIS_BRANCH:-local}
COMMIT_SHORT=${TRAVIS_COMMIT:0:7}

if [ -z "${COMMIT_SHORT}" ]
then
  COMMIT_SHORT=`git rev-parse --short HEAD`
fi

docker build . -t ${APP_NAME}:latest

docker tag ${APP_NAME}:latest ${DOCKER_REPO}:${COMMIT_SHORT}
docker tag ${APP_NAME}:latest ${DOCKER_REPO}:${BRANCH}