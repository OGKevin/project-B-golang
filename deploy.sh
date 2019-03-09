#!/usr/bin/env bash

PROJECT_NAME=project-b-backend

echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin

docker build --build-arg TAG=${TRAVIS_COMMIT} -t ogkevin/${PROJECT_NAME}:${TRAVIS_COMMIT} -f ./cmd/http/Dockerfile .

docker push ogkevin/${PROJECT_NAME}:${TRAVIS_COMMIT}
