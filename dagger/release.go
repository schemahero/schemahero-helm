package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dagger/schemahero-helm/internal/dagger"
)

func (r *SchemaheroHelm) Release(
	ctx context.Context,

	// +defaultPath="./"
	source *dagger.Directory,

	version string,

	// +default=false
	snapshot bool,

	// +default=false
	clean bool,

	onePasswordServiceAccountProduction *dagger.Secret,

	githubToken *dagger.Secret,
) error {
	schemaheroVersion, err := getSchemaheroVersion(ctx, githubToken, source)
	if err != nil {
		return err
	}

	latestVersion, newVersion, err := processVersion(ctx, version, githubToken)
	if err != nil {
		return fmt.Errorf("processing version: %w", err)
	}

	fmt.Printf("Releasing %s -> %s\n", latestVersion, newVersion)

	if err := pushTag(ctx, source, githubToken, newVersion); err != nil {
		return err
	}

	chart := generateChart(ctx, source, newVersion, schemaheroVersion)

	// publish the chart to the ghcr oci registry
	if err := publishChart(ctx, chart, githubToken, newVersion); err != nil {
		return err
	}

	return nil
}

func publishChart(ctx context.Context, chart *dagger.File, githubToken *dagger.Secret, version string) error {
	githubTokenPlaintext, err := githubToken.Plaintext(ctx)
	if err != nil {
		return err
	}

	name, err := chart.Name(ctx)
	if err != nil {
		return err
	}

	// use the helm release container
	container := dag.Container(dagger.ContainerOpts{
		Platform: dagger.Platform("linux/amd64"),
	}).From("alpine/helm:latest").
		WithFile(name, chart).
		WithExec([]string{"helm", "registry", "login", "--username", "schemahero", "--password", githubTokenPlaintext, "ghcr.io/schemahero"}).
		WithExec([]string{"helm", "push", name, "oci://ghcr.io/schemahero/helm"})
	_, err = container.Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}

func generateChart(ctx context.Context, source *dagger.Directory, version string, schemaheroVersion string) *dagger.File {
	// replace the version and appVersion in the helm chart chart.yaml
	chartYAML := dag.Container(dagger.ContainerOpts{
		Platform: dagger.Platform("linux/amd64"),
	}).From("ubuntu:latest").
		WithFile("/chart.yaml", source.File("/Chart.yaml")).
		WithExec([]string{"sed", "-i", "s/version: \"0.0.0\"/version: \"" + version + "\"/g", "/chart.yaml"}).
		WithExec([]string{"sed", "-i", "s/appVersion: \"0.0.0\"/appVersion: \"" + schemaheroVersion + "\"/g", "/chart.yaml"}).
		File("/chart.yaml")
	chartYAMLContents, err := chartYAML.Contents(ctx)
	if err != nil {
		panic(err)
	}

	// generate the helm chart
	chartContainer := dag.Container(dagger.ContainerOpts{
		Platform: dagger.Platform("linux/amd64"),
	}).From("alpine/helm:latest").
		WithDirectory("/chart", source.Directory("/")).
		WithNewFile("/chart/Chart.yaml", chartYAMLContents).
		WithWorkdir("/").
		WithExec([]string{"helm", "package", "./chart", "--destination", "/"})

	chart := chartContainer.
		File(fmt.Sprintf("/schemahero-%s.tgz", version))
	return chart
}

func pushTag(ctx context.Context, source *dagger.Directory, githubToken *dagger.Secret, newVersion string) error {
	err := tryPushTag(ctx, source, githubToken, newVersion)
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "tag already exists in the remote") {
		err = deleteTag(ctx, source, githubToken, newVersion)
		if err != nil {
			return err
		}

		return tryPushTag(ctx, source, githubToken, newVersion)
	}

	return nil
}

func tryPushTag(ctx context.Context, source *dagger.Directory, githubToken *dagger.Secret, newVersion string) error {
	// push the tag
	githubTokenPlaintext, err := githubToken.Plaintext(ctx)
	if err != nil {
		return err
	}
	tagContainer := dag.Container().
		From("alpine/git:latest").
		WithMountedDirectory("/go/src/github.com/schemahero/schemahero-helm", source).
		WithWorkdir("/go/src/github.com/schemahero/schemahero-helm").
		WithExec([]string{"git", "remote", "add", "tag", fmt.Sprintf("https://%s@github.com/schemahero/schemahero-helm.git", githubTokenPlaintext)}).
		WithExec([]string{"git", "tag", newVersion}).
		WithExec([]string{"git", "push", "tag", newVersion})
	_, err = tagContainer.Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}

func deleteTag(ctx context.Context, source *dagger.Directory, githubToken *dagger.Secret, tag string) error {
	fmt.Printf("Deleting tag %s\n", tag)
	// push the tag
	githubTokenPlaintext, err := githubToken.Plaintext(ctx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.github.com/repos/schemahero/schemahero-helm/git/refs/tags/%s", tag), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", githubTokenPlaintext))
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	time.Sleep(time.Second * 10)
	return nil
}
