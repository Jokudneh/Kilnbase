# Kilnbase

[![Go 1.22+](https://img.shields.io/badge/go-%3E%3D1.22-00ADD8.svg)](go.mod)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

**Image build tool for Kubernetes:** reads a config, resolves every image reference in it and
writes the digest-pinned result back.


[![OpenSSF Best Practices](https://bestpractices.coreinfrastructure.org/projects/7746/badge)](https://bestpractices.coreinfrastructure.org/projects/7746)

# kilnbase

- Website: https://github.com/Jokudneh/Kilnbase
- Slack: [#kilnbase in Kubernetes slack](https://slack.kubernetes.io)
- [Docs](docs/README.md) with topics about building, packaging images, blog posts, etc.
- Install: Grab prebuilt binaries from the Releases page or [Homebrew Kilnbase tap](https://github.com/carvel-dev/homebrew)
- Backlog: [See what we're up to](https://github.com/orgs/kilnbase-dev/projects/1/views/1?filterQuery=repo%3A%22kilnbase-dev%2Fkilnbase%22).

`kilnbase` (pronounced: `kei·bild`) seamlessly incorporates image building and image pushing into your development and deployment workflows.

Features:

- Orchestrates image builds (delegates to tools like Docker, pack, kubectl-buildkit) and registry pushes
- Works with local Docker daemon and remote registries, for development and production cases
- Records metadata about image sources in annotation on Kubernetes resources (see examples below)
- Resolves image references to their digest form (*immutable*) ([details](https://github.com/Jokudneh/Kilnbase/docs/latest/#why-digest-references))
- Provides a way to transport set of images in a single tarball between registries
  - maintaining exactly same digests, hence guaranteeing exactly same content
- Not specific to Kubernetes, but works really well with Kubernetes configuration files  

!

See [building and deploying simple Go application to Kubernetes example](https://github.com/carvel-dev/simple-app-on-kubernetes#step-3-building-container-images-locally) that uses kilnbase.

### Join the Community and Make Kilnbase Better
Kilnbase is better because of our contributors and maintainers. It is because of you that we can bring great software to the community.
Please join us during our online community meetings. Details can be found on our [Kilnbase website](https://carvel.dev/community/).

You can chat with us on Kubernetes Slack in the #kilnbase channel and follow us on Twitter at @kilnbase_dev.

Check out which organizations are using and contributing to Kilnbase: [Adopter's list](https://github.com/carvel-dev/kilnbase/blob/master/ADOPTERS.md)

# Development

Consult [docs/dev.md](docs/dev.md) for build instructions, code structure details.