// Package github provides a GitHub release store.
package github

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/tj/go-update"
)

// Store is the store implementation.
type Store struct {
	Owner   string
	Repo    string
	Version string
}

/*
GetRelease returns the specified release or ErrNotFound.

It tries to find a GitHub release with a tag matching the given version
with "v" prefixed.

If this first step fails, it tries to find a GitHub release with a tag matching
exactly with the given version.
*/
func (s *Store) GetRelease(version string) (*update.Release, error) {
	r, err := getGithubReleaseByTag("v" + version)

	if err != nil {
		if _, ok := err.(*update.ErrNotFound); ok {
			return getGithubReleaseByTag(version)
		}

		return nil, err
	}

	return r, nil
}

// LatestReleases returns releases newer than Version, or nil.
func (s *Store) LatestReleases() (latest []*update.Release, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gh := github.NewClient(nil)

	releases, _, err := gh.Repositories.ListReleases(ctx, s.Owner, s.Repo, nil)
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		tag := r.GetTagName()

		if tag == s.Version || "v"+s.Version == tag {
			break
		}

		latest = append(latest, githubRelease(r))
	}

	return
}

// getGithubReleaseByTag returns the specified release or ErrNotFound.
func getGithubReleaseByTag(tag string) (*update.Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gh := github.NewClient(nil)

	r, res, err := gh.Repositories.GetReleaseByTag(ctx, s.Owner, s.Repo, tag)

	if res.StatusCode == 404 {
		return nil, update.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return githubRelease(r), nil
}

// githubRelease returns a Release.
func githubRelease(r *github.RepositoryRelease) *update.Release {
	out := &update.Release{
		Version:     r.GetTagName(),
		Notes:       r.GetBody(),
		PublishedAt: r.GetPublishedAt().Time,
		URL:         r.GetURL(),
	}

	for _, a := range r.Assets {
		out.Assets = append(out.Assets, &update.Asset{
			Name:      a.GetName(),
			Size:      a.GetSize(),
			URL:       a.GetBrowserDownloadURL(),
			Downloads: a.GetDownloadCount(),
		})
	}

	return out
}
