RELEASE_VERSION=0.1.1

all: test policy

run:
	go run *.go

policy:
	PROJECT_NUMBER="$(gcloud projects describe ${PROJECT_ID} --format='get(projectNumber)')"
	gcloud projects add-iam-policy-binding ${PROJECT_NUMBER} \
    	--member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
    	--role=roles/container.developer

image:
	gcloud builds submit \
		--project=$(PROJECT_ID) \
		--tag gcr.io/$(PROJECT_ID)/maxprime .

sample-image:
	gcloud builds submit \
		--project=knative-samples \
		--tag gcr.io/knative-samples/maxprime .

trigger:
	# TODO: Implement on-demand build
	# https://github.com/kelseyhightower/pipeline/blob/master/labs/build-triggers.md
	# https://cloudbuild.googleapis.com/v1/projects/$%7BPROJECT_ID%7D/triggers

deploy:
	kubectl apply -f deployments/service.yaml -n demo

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline