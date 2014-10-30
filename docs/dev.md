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
