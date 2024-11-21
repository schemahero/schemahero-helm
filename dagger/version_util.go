package main

import (
	"context"
	"dagger/schemahero-helm/internal/dagger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver"
)

func getSchemaheroVersion(ctx context.Context, githubToken *dagger.Secret, source *dagger.Directory) (string, error) {
	latestVersion, err := getLatestVersion(ctx, "schemahero", "schemahero", githubToken)
	if err != nil {
		return "", err
	}
	return latestVersion, nil
}

func processVersion(ctx context.Context, requestedVersion string, githubToken *dagger.Secret) (string, string, error) {
	latestVersion, err := getLatestVersion(ctx, "schemahero", "schemahero-helm", githubToken)
	if err != nil {
		return "", "", err
	}

	parsedLatestVersion, err := semver.NewVersion(latestVersion)
	if err != nil {
		return "", "", err
	}

	switch requestedVersion {
	case "major":
		return latestVersion, fmt.Sprintf("%d.0.0", parsedLatestVersion.Major()+1), nil
	case "minor":
		return latestVersion, fmt.Sprintf("%d.%d.0", parsedLatestVersion.Major(), parsedLatestVersion.Minor()+1), nil
	case "patch":
		return latestVersion, fmt.Sprintf("%d.%d.%d", parsedLatestVersion.Major(), parsedLatestVersion.Minor(), parsedLatestVersion.Patch()+1), nil
	default:
		return latestVersion, requestedVersion, nil
	}
}

func getLatestVersion(ctx context.Context, owner string, repo string, githubToken *dagger.Secret) (string, error) {
	githubTokenPlain, err := githubToken.Plaintext(ctx)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprint("https://api.github.com/repos/%s/%s/releases/latest", owner, repo), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", githubTokenPlain))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}
