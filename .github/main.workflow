workflow "Publish container to Docker" {
  on       = "push"

  resolves = [
    "Push"
  ]
}

action "Test" {
  uses = "./.github/actions/golang"
  args = "test"
}

action "Build" {
  uses    = "./.github/actions/docker"

  needs   = [
    "Test"
  ]

  secrets = [
    "DOCKER_IMAGE"
  ]

  args    = [
    "build",
    "Dockerfile"
  ]
}

action "Login" {
  uses    = "actions/docker/login@master"

  needs   = [
    "Build"
  ]

  secrets = [
    "DOCKER_USERNAME",
    "DOCKER_PASSWORD"
  ]
}

action "Push" {
  uses    = "./.github/actions/docker"
  args    = "push"

  needs   = [
    "Login"
  ]

  secrets = [
    "DOCKER_IMAGE"
  ]
}
