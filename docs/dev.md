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
- [pkg/kilnbase/search](https://github.com/kilnbase-dev/kilnbase/tree/develop/pkg/kilnbase/search) implements YAML node searcher that finds image URLs
- [test/e2e](https://github.com/kilnbase-dev/kilnbase/tree/develop/test/e2e) includes e2e tests that can run against Docker registry.

## How to set up an insecure-only registry
It is occasionally necessary to test against insecure-only registries in order to ensure the commands can run against http endpoints.
Here's how we set one up on Google's Kubernetes Engine:


* Create a GKE cluster through the web GUI
* Install gcloud - https://cloud.google.com/kubernetes-engine/docs/quickstart#local-shell
* Target your cluster `gcloud container clusters get-credentials <cluster-name>`
* Create the registry deployment file, reg.yml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: registry
  labels:
    app: registry
spec:
  replicas: 1
  selector:
    matchLabels:
      app: registry
  template:
    metadata:
      labels:
        app: registry
    spec:
      containers:
      - name: registry
        image: registry
        ports:
        - containerPort: 8080
        env:
        - name: REGISTRY_PORT
          value: "8080"
```

* Apply the deployment with `kubectl apply -f reg.yml`

* Expose it on port 5000 with
`kubectl expose deployment registry --name=registry --type=LoadBalancer --port 80 --target-port 5000 -oyaml`

* Wait until it gets an external IP (~30 seconds)

To verify that the registry is up and running:\
In [GKE](https://console.cloud.google.com/), navigate to Kubernetes Engine, then 'Services & Ingress'\
Confirm that there is an 'External load balancer' named 'registry' with an 'External Endpoint' and a link that returns a 200.\
At the bottom of the details view (click on registry), see that port forwarding is set up with a 'target port' of 5000.

That's all!  We should now be able to run regular commands against it:
`kilnbase relocate -f test/e2e/assets/simple-app -r <External IP>:80/<Name> --registry-insecure`
