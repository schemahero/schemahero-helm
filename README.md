# Helm chart for SchemaHero

## Installation

The chart is published to
[GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
and can be installed only with [Helm 3](https://helm.sh/docs).

Helm 3 CLI requires enabling [OCI support](https://helm.sh/docs/topics/registries/#enabling-oci-support).

```sh
export HELM_EXPERIMENTAL_OCI=1
```

First, find the appropriate version from
[available releases](https://github.com/eshepelyuk/schemahero-helm/releases) list.

Then, download the chart with choosen `${VERSION}` to a local directory.

```sh
helm pull oci://ghcr.io/eshepelyuk/helm/schemahero --version ${VERSION}
```

The command above will download a file named `schemahero-${VERSION}.tgz`.

After that, chart is ready to be installed (or upgraded),
we suggest to install it into a dedicated namespace.

```sh
helm upgrade -i --wait --create-namespace -n schemahero schemahero schemahero-${VERSION}.tgz
```

## Configuration

Chart parameters can be configured via [Helm values files](https://helm.sh/docs/chart_template_guide/values_files/).

Check out
[values schema](https://artifacthub.io/)
for the available configuration options.

## Create new release

1. Navigate to [Release Workflow](https://github.com/eshepelyuk/schemahero-helm/actions/workflows/release.yaml)
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

