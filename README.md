# maxprime

UI and REST service demo calculating highest prime number up to the number specified by the user. Good for infrastructure scaling demo

[![Click to run on Cloud Run](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/mchmarny/maxprime.git)

> Note, this app is also used to demonstrate GitOps on Knative using Cloud Build. See the GitOps section below

## Demo

Live version of this app is available at
https://maxprime.demo.knative.tech

### Request

To calculate the highest prime number in UI

```
curl http://localhost:8080
```

You can also pass max number using REST service

```
curl http://localhost:8080/prime/5000000
```

### Response

`prime.val` in the response has the highest prime number bellow passed in argument.

```
{
    "id": "0fe3fa30-627e-46cc-be09-a9de62252ff1",
    "ts": "2019-03-24 13:54:14.917105 +0000 UTC",
    "dur": "724.848232ms",
    "rel": "v0.0.1",
    "prime": {
        "max": 50000340,
        "val": 50000329
    }
}
```

## GitOps

Simple setup to automate Knative deployments using Git and Cloud Build

As a developer, you write code and commit it to a git repo. You also hopefully run tests on that code for each commit. Assuming your application passes all the tests, you may want to deploy it to Knative cluster. You can do it form your workstation by using any one of the Knative CLIs (e.g. gcloud, knctl, tm etc.).

In this demo however we are going to demonstrate deploying directly from git repository. This means that you as a developer do not need install anything on your machine other than the standard git tooling. Here is the outline:

* Create a release tag on the commit you want to deploy in git
* Cloud Build then:
  * Tests (again)
  * Builds and tags image
  * Pushes that image to repository
  * Creates Knative service manifest
  * Applies that manifest to designated Knative cluster

> As an add-on, we are also going to send mobile notification with build status using [knative-build-status-notifs](https://github.com/mchmarny/knative-build-status-notifs)

### Setup

You will be using a number of GCP APIs, so you can enable them all now: 
```
gcloud services enable \
    run.googleapis.com \
    cloudbuild.googleapis.com
```

You will have to [configure git trigger](https://console.cloud.google.com/cloud-build/triggers/add) in Cloud Build first. There doesn't seem to be a way to do this using `gcloud`.

> Trigger type: Tag
> Tag (regex): `release-*`
> Build configuration: `deployments/cloudbuild.yaml`

Then setup IAM policy binding to allow Cloud Builder deploy build image to your cluster

```shell
PROJECT_NUMBER="$(gcloud projects describe ${PROJECT_ID} --format='get(projectNumber)')"
gcloud projects add-iam-policy-binding ${PROJECT_NUMBER} \
    --member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
    --role=roles/container.developer
```

You will also need to change the `deployments/cloudbuild.yaml` file to name your app and cluster information. 


Finally submit the Cloud Build configuration

```shell
gcloud builds submit --config deployments/cloudbuild.yaml
```

## Deployment

To build and deploy specific commit from git, tag it and publish the tags. We also are going to print the last few tags so we can see the exact commit hash.

```shell
git tag "release-v${RELEASE_VERSION}"
git push origin "release-v${RELEASE_VERSION}"
git log --oneline
```

### Logs

You can monitor progress of your build but first finding its id

```shell
gcloud builds list
```

And then describing it

```shell
gcloud builds describe BUILD_ID
```

You can always also navigate to the [Build History](https://console.cloud.google.com/cloud-build/builds) screen in UI and see it there.

## Load


First forward Knative Monitoring (Grafan) to a local port (`3000`)

```shell
kubectl port-forward -n knative-monitoring \
    $(kubectl get pods -n knative-monitoring --selector=app=grafana \
    --output=jsonpath="{.items..metadata.name}") 3000
```

Now you can view the Knative revision dashboard [here](http://localhost:3000/d/im_gFbWik/knative-serving-revision-http-requests?refresh=3s&orgId=1&var-namespace=demo&var-configuration=maxprime&var-revision=All)

You can also view the `maxprime` pods beign scaled here

```shell
watch kubectl get pods -n demo -l serving.knative.dev/service=maxprime
```

And finally, to generate some load, run the included load generator. `--count` is the number of concurrent client (threads), `--prime` is the max number to calculate prime for, and `--url` is the target URL where `maxprime` demo is deployed (no final slash).

```
go run load/main.go --count 300 --prime 99999 \
    --url https://maxprime.demo.knative.tech
```
