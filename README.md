# registry-operator

Operator for CNCF Distribution Registry 📦

[![GitHub License](https://img.shields.io/github/license/registry-operator/registry-operator)][license]
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
[![GitHub issues](https://img.shields.io/github/issues/registry-operator/registry-operator)](https://github.com/registry-operator/registry-operator/issues)
[![GitHub release](https://img.shields.io/github/release/registry-operator/registry-operator)](https://GitHub.com/registry-operator/registry-operator/releases/)
[![Go Reference](https://pkg.go.dev/badge/github.com/registry-operator/registry-operator)](https://pkg.go.dev/github.com/registry-operator/registry-operator)

Join the community to discuss ongoing operator development and usage in [#registry-operator](https://cloud-native.slack.com/archives/C06P7RC8857) channel on the CNCF Slack and see the [self generated docs](https://pkg.go.dev/github.com/registry-operator/registry-operator) for more usage details.

[![Slack Channel](https://img.shields.io/badge/Slack-CNCF-4A154B?logo=slack)](https://cloud-native.slack.com/archives/C06P7RC8857)

The following quality gates are ensuring the quality and maintainability of the codebase.

[![Go Report Card](https://goreportcard.com/badge/github.com/registry-operator/registry-operator)](https://goreportcard.com/report/github.com/registry-operator/registry-operator)
[![codecov](https://codecov.io/gh/registry-operator/registry-operator/graph/badge.svg?token=TDD92A90UE)](https://codecov.io/gh/registry-operator/registry-operator)

## Table of Contents

- [registry-operator](#registry-operator)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [Getting Started](#getting-started)
  - [Contributing](#contributing)
  - [License](#license)

## Overview

The `registry-operator` is a powerful tool designed to manage and operate CNCF Distribution Registry instances. It streamlines the deployment, scaling, and management of container image registries, providing a seamless experience for developers and DevOps teams.

## Features

- **Automated Deployment**: Deploy CNCF Distribution Registry instances with ease using declarative configuration.
- **Scalability**: Scale your registry horizontally to handle increased loads and ensure high availability.
- **Monitoring and Alerting**: Gain insights into the health and performance of your registry with built-in monitoring and alerting features.
- **Resource Optimization**: Optimize resource utilization by dynamically adjusting resources based on workload demands.

## Getting Started

To start using `registry-operator`, follow these steps:

1. **Installation**: Install the operator in your Kubernetes cluster using Helm or Kubernetes manifests.
2. **Configuration**: Customize the configuration to match your environment and requirements.
3. **Deployment**: Deploy CNCF Distribution Registry instances using the operator.
4. **Management**: Manage and monitor your registry instances using the provided tools and APIs.

For detailed instructions, refer to the [documentation][documentation].

## Contributing

We welcome contributions from the community! If you're interested in contributing to `registry-operator`, please read our [Contribution Guidelines](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) before getting started.

## License

`registry-operator` is licensed under the [Apache-2.0][license].

<!-- Resources -->

[documentation]: https://registry-operator.github.io/docs
[license]: https://github.com/registry-operator/registry-operator/blob/main/LICENSE
