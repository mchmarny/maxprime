RELEASE_VERSION=0.4.1
APP_NS?=demo
CLUSTER_NAME?=demo
IMAGE_OWNER ?=$(shell git config --get user.username)

.PHONY: run policy deploy image ghimage kimage trigger apply tag

all: test

tidy:
	go mod tidy
	go mod vendor

run:
	go run *.go

policy:
	gcloud projects add-iam-policy-binding ${PROJECT_NUMBER} \
    	--member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
    	--role=roles/container.developer

deploy: tidy
	gcloud builds submit \
		--project=$(PROJECT_ID) \
		--config=deployments/cloudbuild.yaml \
		--substitutions=_APP_NAME=maxprime,_APP_NS=$(APP_NS),_CLUSTER_NAME=kn07,_CLUSTER_ZONE=$(CLUSTER_ZONE),SHORT_SHA=$(COMMIT_SHA) \
		.

image: tidy
	gcloud builds submit \
		--project=cloudylabs \
		--tag "gcr.io/cloudylabs/maxprime:${RELEASE_VERSION}" .

ghimage: tidy ## Builds and publish image 
	docker build -t "ghcr.io/$(IMAGE_OWNER)/maxprime:$(RELEASE_VERSION)" .
	docker push "ghcr.io/$(IMAGE_OWNER)/maxprime:$(RELEASE_VERSION)"

apply:
	kubectl apply -f deployments/service.yaml -n demo

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline