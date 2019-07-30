RELEASE_VERSION=0.1.1

all: test policy

run:
	go run *.go

policy:
	PROJECT_NUMBER="$(gcloud projects describe ${PROJECT_ID} --format='get(projectNumber)')"
	gcloud projects add-iam-policy-binding ${PROJECT_NUMBER} \
    	--member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
    	--role=roles/container.developer

deploy:
	gcloud builds submit \
		--project=$(PROJECT_ID) \
		--config=deployments/cloudbuild.yaml \
		--substitutions=_APP_NAME=maxprime,_APP_NS=demo,_CLUSTER_NAME=kn07,_CLUSTER_ZONE=us-west1-c \
		.

image:
	gcloud builds submit \
		--project=cloudylabs-public \
		--tag gcr.io/cloudylabs-public/maxprime .

kimage:
	gcloud builds submit \
		--project=knative-samples \
		--tag gcr.io/knative-samples/maxprime .

trigger:
	# TODO: Implement on-demand build
	# https://github.com/kelseyhightower/pipeline/blob/master/labs/build-triggers.md
	# https://cloudbuild.googleapis.com/v1/projects/$%7BPROJECT_ID%7D/triggers

apply:
	kubectl apply -f deployments/service.yaml -n demo

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline