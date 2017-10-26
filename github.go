package update

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// Github store.
type Github struct {
	Owner   string
	Repo    string
	Version string
}

// LatestReleases returns releases newer than Version, or nil.
func (s *Github) LatestReleases() (latest []*Release, err error) {
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

		latest = append(latest, toRelease(r))
	}

	return
}

// toRelease returns a Release.
func toRelease(r *github.RepositoryRelease) *Release {
	out := &Release{
		Version:     r.GetTagName(),
		Notes:       r.GetBody(),
		PublishedAt: r.GetPublishedAt().Time,
		URL:         r.GetURL(),
	}

	for _, a := range r.Assets {
		out.Assets = append(out.Assets, &Asset{
			Name:      a.GetName(),
			Size:      a.GetSize(),
			URL:       a.GetBrowserDownloadURL(),
			Downloads: a.GetDownloadCount(),
		})
	}

	return out
}
