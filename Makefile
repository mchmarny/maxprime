RELEASE_VERSION=0.2.3
PROJECT_NUMBER=$(shell gcloud projects describe ${PROJECT_ID} --format='get(projectNumber)')
COMMIT_SHA=$(shell git rev-parse HEAD)
APP_NS?=demo
CLUSTER_NAME?=kn07
CLUSTER_ZONE?=us-west1-c
DOCKER_HUB_USER?=mchmarny

.PHONY: run policy deploy image kimage trigger apply tag

all: test

mod:
	go mod tidy
	go mod vendor

run:
	go run *.go

policy:
	gcloud projects add-iam-policy-binding ${PROJECT_NUMBER} \
    	--member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
    	--role=roles/container.developer

deploy: mod
	gcloud builds submit \
		--project=$(PROJECT_ID) \
		--config=deployments/cloudbuild.yaml \
		--substitutions=_APP_NAME=maxprime,_APP_NS=$(APP_NS),_CLUSTER_NAME=kn07,_CLUSTER_ZONE=$(CLUSTER_ZONE),SHORT_SHA=$(COMMIT_SHA) \
		.

image: mod
	gcloud builds submit \
		--project=cloudylabs \
		--tag "gcr.io/cloudylabs/maxprime:${RELEASE_VERSION}" .

publicimage: mod
	gcloud builds submit \
		--project=cloudylabs-public \
		--tag "gcr.io/cloudylabs-public/maxprime:${RELEASE_VERSION}" .

dockerpush: mod
	docker build -t "${DOCKER_HUB_USER}/maxprime:${RELEASE_VERSION}" .
	docker push "${DOCKER_HUB_USER}/maxprime:${RELEASE_VERSION}"

apply:
	kubectl apply -f deployments/service.yaml -n demo

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline