
[![OpenSSF Best Practices](https://bestpractices.coreinfrastructure.org/projects/7746/badge)](https://bestpractices.coreinfrastructure.org/projects/7746)

# kilnbase

- Website: https://github.com/Jokudneh/Kilnbase
- Slack: [#kilnbase in Kubernetes slack](https://slack.kubernetes.io)
- [Docs](docs/README.md) with topics about building, packaging images, blog posts, etc.
- Install: Grab prebuilt binaries from the Releases page or [Homebrew Kilnbase tap](https://github.com/kilnbase-dev/homebrew)
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
