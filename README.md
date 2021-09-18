
# Helm chart for SchemaHero

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/schemahero/schemahero-helm/CI?label=CI%2FCD&style=for-the-badge)
[![Current](https://img.shields.io/github/v/tag/schemahero/schemahero-helm?logo=github&sort=semver&style=for-the-badge&label=current)](https://github.com/schemahero/schemahero-helm/releases/latest)
[![Apache 2.0](https://img.shields.io/github/license/schemahero/schemahero-helm?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)

## Installation

The chart is published to
[GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
and can be installed only with [Helm 3](https://helm.sh/docs).

!!! Helm version must be >= 3.7.0

First, enable [OCI support](https://helm.sh/docs/topics/registries/#enabling-oci-support).

```sh
export HELM_EXPERIMENTAL_OCI=1
```

Choose appropriate version from
[available releases](https://github.com/schemahero/schemahero-helm/releases) list.

It's recommended to install the chart into a dedicated namespace.

```sh
helm upgrade -i --wait --create-namespace -n schemahero schemahero \
  oci://ghcr.io/schemahero/helm/schemahero --version ${VERSION}
```

## Configuration

Chart parameters can be configured via [Helm values files](https://helm.sh/docs/chart_template_guide/values_files/).

Check out
[values schema](https://artifacthub.io/)
for the available configuration options.

## Create new release

1. Navigate to [Release Workflow](https://github.com/schemahero/schemahero-helm/actions/workflows/release.yaml)
in `Actions` section.
1. Click `Run workflow` button.
1. Wait until job is finished.

Release job performs following actions.

1. Generates changelog from previous release tag.
1. Calculates a new release version from the changelog.
1. Creates new GitHub release with git tag containing new release number.
1. Generates Helm chart with assigned version.
1. Publishes Helm chart to GitHub Packages repository.

Release pipeline is implemented with
[semantic-release](https://github.com/semantic-release/semantic-release) tool.

