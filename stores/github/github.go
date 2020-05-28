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

type version []int64

// GetRelease returns the specified release or ErrNotFound.
func (s *Store) GetRelease(version string) (*update.Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gh := github.NewClient(nil)

	r, res, err := gh.Repositories.GetReleaseByTag(ctx, s.Owner, s.Repo, "v"+version)

	if res.StatusCode == 404 {
		return nil, update.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return githubRelease(r), nil
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

	vsn := parseVersion(s.Version)
	for _, r := range releases {
		tag := parseVersion(r.GetTagName())

		if vsn.biggerEquals(tag) {
			break
		}

		latest = append(latest, githubRelease(r))
	}

	return
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

func parseVersion(vsn string) version {
	vsn = strings.TrimLeft(vsn, "v")
	vsn = strings.ReplaceAll(vsn, "-", ".")
	parts := strings.Split(vsn, ".")
	ret := make(version, len(parts))
	for n, p := range parts {
		i, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			i = 0
		}
		ret[n] = i
	}
	return ret
}

func (vsn version) at(pos int) int64 {
	if len(vsn) <= pos {
		return 0
	}
	return vsn[pos]
}

func (vsn version) biggerEquals(other version) bool {
	maximum := len(vsn)
	if len(other) > maximum {
		maximum = len(other)
	}
	for i := 0; i <= maximum; i++ {
		if vsn.at(i) == other.at(i) {
			continue
		}
		return vsn.at(i) > other.at(i)
	}
	return true
}
