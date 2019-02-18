#!/usr/bin/env bash

PROJECT_NAME=project-b-backend

docker build --build-arg TAG=${TRAVIS_COMMIT} -t ${PROJECT_NAME} -f ./cmd/http/Dockerfile .

docker tag ${PROJECT_NAME} ${PROJECT_NAME}:latest
docker tag $1:latest ogkevin/${PROJECT_NAME}:latest

docker tag ${PROJECT_NAME} ${PROJECT_NAME}:${TRAVIS_COMMIT}
docker tag ${PROJECT_NAME}:${TRAVIS_COMMIT} ogkevin/${PROJECT_NAME}:${TRAVIS_COMMIT}
