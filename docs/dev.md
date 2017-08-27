# Development

## Prerequisites

- [minikube](https://minikube.sigs.k8s.io/docs/)
- [docker](https://www.docker.com/)

## Run Unit tests
```bash
# Run all tests
./hack/test.sh
# or run single test
./hack/test.sh -run TestLogger
```

## Run E2E tests against minikube registry
Install virtualbox from https://www.virtualbox.org/wiki/Downloads. After installing the first time, you may need to reboot.
```bash
# Bootstrap k8s cluster and enable docker registry
# X.X.X.X must be replaced with your subnetmask of "minikube ip"
minikube start --driver=virtualbox --insecure-registry=X.X.X.X/16 --host-only-cidr "192.168.59.1/24"
# Build kilnbase binary for testing
./hack/build.sh
# Make your env aware of the docker registry
eval $(minikube docker-env)
# Run all tests
./hack/test-all-locally.sh
# or run single test
./hack/test-all-locally.sh -run TestDockerBuildSuccessful
```

## Source Code Structure

For those interested in extending and improving kilnbase, below is a quick reference on the structure of the source code:

- [.github/workflows/test-gh.yml](https://github.com/kilnbase-dev/kilnbase/blob/develop/.github/workflows/test-gh.yml) is a Github Action that runs build and unit tests when commits are pushed
- [hack](https://github.com/kilnbase-dev/kilnbase/tree/develop/hack) has build and test scripts
- [cmd/kilnbase](https://github.com/kilnbase-dev/kilnbase/blob/develop/cmd/kilnbase) is the entry package for main kilnbase binary
- [pkg/kilnbase/cmd](https://github.com/kilnbase-dev/kilnbase/tree/develop/pkg/kilnbase/cmd) includes all kilnbase CLI commands (kilnbase.go is root command)
- [pkg/kilnbase/config](https://github.com/kilnbase-dev/kilnbase/tree/develop/pkg/kilnbase/config) describes kilnbase configuration resources such as Config, Sources, etc.
- [pkg/kilnbase/resources](https://github.com/kilnbase-dev/kilnbase/tree/develop/pkg/kilnbase/resources) allows to parse YAML files into Resource objects 
- [pkg/kilnbase/image](https://github.com/kilnbase-dev/kilnbase/tree/develop/pkg/kilnbase/image) contains set of classes that know how to "transform" image URLs (build it, resolve it to digest, tag it)
- [pkg/kilnbase/registry](https://github.com/kilnbase-dev/kilnbase/tree/develop/pkg/kilnbase/registry) provides a simplified registry API client
