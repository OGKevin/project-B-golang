#!/usr/bin/env bash

PROJECT_NAME=project-b-backend

echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin

docker build --build-arg TAG=${TRAVIS_COMMIT} -t ogkevin/${PROJECT_NAME}:${TRAVIS_COMMIT} -f ./cmd/http/Dockerfile .

docker push ogkevin/${PROJECT_NAME}:${TRAVIS_COMMIT}

curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl

chmod +x ./do-linux
./do-linux kube_config

export KUBECONFIG=$(pwd)/do-kube.yaml

kubectl config use-context do-ams3-k8s-1-13-2-do-1-ams3-ogkevin

kubectl -n project-b set image deployment --record=true -l project=project-b project-b-backed=ogkevin/project-b-backend:${TRAVIS_COMMIT}
kubectl rollout status deployment project-b-backed

curl -D - -S -XPOST -H "Github-Webhook-Secret: ${GITHUB_WEBHOOK_SECRET}" -H "Content-type: application/json" -d "{\"sha\": \"${TRAVIS_COMMIT}\", \"owner\":\"OGKevin\", \"repo\":\"project-b-golang\"}" 'https://ogkevin.nl/api/github/release'
