## Hipster Shop

Repo: [https://github.com/GoogleCloudPlatform/microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo)

```bash
$ git clone https://github.com/GoogleCloudPlatform/microservices-demo
$ kapp deploy -a md -f <(kilnbase -f microservices-demo/kubernetes-manifests -f ../examples/gcp-microservices-demo/prebuilt.yml)
```

With your own application one would expect that `prebuilt.yml` configuration would be part of that repository (instead of similar to this example, where `prebuilt.yml` is in a separate repo).
