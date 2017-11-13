// Package apex provides an Apex release store.
package apex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/go-update"
	"github.com/tj/go/http/request"
)

// Store is the store implementation.
type Store struct {
	URL       string
	Product   string
	Plan      string
	Version   string
	AccessKey string
}

// Release model.
type Release struct {
	Version   string    `json:"version"`
	Notes     string    `json:"notes"`
	Files     []*File   `json:"files"`
	CreatedAt time.Time `json:"created_at"`
}

// File model.
type File struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// LatestReleases returns releases newer than Version, or nil.
func (s *Store) LatestReleases() (latest []*update.Release, err error) {
	url := fmt.Sprintf("%s/%s/%s", s.URL, s.Product, s.Plan)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	req.Header.Set("Authorization", "Bearer "+s.AccessKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "requesting")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, request.Error(res.StatusCode)
	}

	var releases []*Release

	if err := json.NewDecoder(res.Body).Decode(&releases); err != nil {
		return nil, errors.Wrap(err, "unmarshaling")
	}

	for _, r := range releases {
		if r.Version == s.Version {
			break
		}

		latest = append(latest, toRelease(r))
	}

	return
}

// toRelease returns a Release.
func toRelease(r *Release) *update.Release {
	out := &update.Release{
		Version:     r.Version,
		Notes:       r.Notes,
		PublishedAt: r.CreatedAt,
	}

	for _, f := range r.Files {
		out.Assets = append(out.Assets, &update.Asset{
			Name: f.Name,
			Size: int(f.Size),
			URL:  f.URL,
		})
	}

	return out
}
