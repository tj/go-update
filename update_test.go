package update

import (
	"testing"

	"github.com/apex/log"
	"github.com/tj/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// TODO: mock github, this will obviously break

func TestProject_LatestReleases(t *testing.T) {
	t.Run("when old", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			Command: "polls",
			Owner:   "tj",
			Repo:    "gh-polls",
			Version: "0.0.3",
		}

		releases, err := p.LatestReleases()
		assert.NoError(t, err)

		var versions []string

		for _, r := range releases {
			versions = append(versions, r.Version)
		}

		assert.Equal(t, []string{"v0.1.1", "v0.1.0"}, versions)
	})

	t.Run("when new", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			Command: "polls",
			Owner:   "tj",
			Repo:    "gh-polls",
			Version: "0.1.1",
		}

		releases, err := p.LatestReleases()
		assert.NoError(t, err)
		assert.Nil(t, releases, "releases")
	})
}

func TestRelease_FindTarball(t *testing.T) {
	t.Parallel()

	p := &Project{
		Command: "polls",
		Owner:   "tj",
		Repo:    "gh-polls",
		Version: "0.0.3",
	}

	releases, err := p.LatestReleases()
	assert.NoError(t, err)
	assert.NotNil(t, releases, "releases")
	r := releases[0]

	a := r.FindTarball("darwin", "amd64")
	assert.NotNil(t, a, "nil for darwin")
	assert.Equal(t, "gh-polls_0.1.1_darwin_amd64.tar.gz", a.Name)

	a = r.FindTarball("windows", "amd64")
	assert.NotNil(t, a, "nil for windows")
	assert.Equal(t, "gh-polls_0.1.1_windows_amd64.tar.gz", a.Name)

	a = r.FindTarball("sloth", "amd64")
	assert.Nil(t, a)
}

func TestAsset_Download(t *testing.T) {
	t.Parallel()

	p := &Project{
		Command: "polls",
		Owner:   "tj",
		Repo:    "gh-polls",
		Version: "0.0.3",
	}

	releases, err := p.LatestReleases()
	assert.NoError(t, err)
	assert.NotNil(t, releases, "releases")
	r := releases[0]

	a := r.FindTarball("darwin", "amd64")
	assert.NotNil(t, a, "nil for darwin")

	path, err := a.Download()
	assert.NoError(t, err, "download")
	assert.NotEmpty(t, path, "path")
}

func TestProject_Install(t *testing.T) {
	t.Parallel()

	p := &Project{
		Command: "polls",
		Owner:   "tj",
		Repo:    "gh-polls",
		Version: "0.0.3",
	}

	releases, err := p.LatestReleases()
	assert.NoError(t, err)
	assert.NotNil(t, releases, "releases")
	r := releases[0]

	a := r.FindTarball("darwin", "amd64")
	assert.NotNil(t, a, "nil for darwin")

	path, err := a.Download()
	assert.NoError(t, err, "download")
	assert.NotEmpty(t, path, "path")

	err = p.Install(path)
	assert.NoError(t, err, "replace")
}
