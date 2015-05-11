## End-to-end test

### Prerequisites

To run the end to end tests, you must have the following utilities installed on your test runner. The test script checks for the specific version referenced. While a package installer will install the utility, it may install a newer, or different version than the one specified. Be sure to install these versions.

- [pack 0.8.1](https://github.com/buildpacks/pack)
- [kubectl-buildkit 0.1.0](https://github.com/jokudneh/buildkit-cli-for-kubectl/releases/tag/v0.1.0)
- [ko 0.8.0](https://github.com/google/ko/releases/tag/v0.8.0)
- [bazel 4.2.0](https://github.com/bazelbuild/bazel/releases/tag/4.2.0)

### Run End to End Tests

Set environment variables

```bash
export KILNBASE_E2E_DOCKERHUB_USERNAME=joeexample
```

To use a registry other than DockerHub, set the KILNBASE_E2E_DOCKERHUB_HOSTNAME variable. You can specify a local registry or use an anonymous service, such as ttl.sh. You will need to authenticate with non-anonymous registeries for the tests to pass.

```bash
