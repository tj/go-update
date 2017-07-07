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
			Command: "apex",
			Owner:   "apex",
			Repo:    "apex",
			Version: "0.13.1",
		}

		releases, err := p.LatestReleases()
		assert.NoError(t, err)

		var versions []string

		for _, r := range releases {
			versions = append(versions, r.Version)
		}

		assert.Equal(t, []string{"v0.15.0", "v0.14.0"}, versions)
	})

	t.Run("when new", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			Command: "apex",
			Owner:   "apex",
			Repo:    "apex",
			Version: "v0.15.0",
		}

		releases, err := p.LatestReleases()
		assert.NoError(t, err)
		assert.Nil(t, releases, "releases")
	})
}

func TestRelease_Asset(t *testing.T) {
	t.Parallel()

	p := &Project{
		Command: "apex",
		Owner:   "apex",
		Repo:    "apex",
		Version: "0.13.1",
	}

	releases, err := p.LatestReleases()
	assert.NoError(t, err)
	assert.NotNil(t, releases, "releases")
	r := releases[0]

	a := r.Asset("darwin", "386")
	assert.NotNil(t, a, "nil for darwin")
	assert.Equal(t, "apex_darwin_386", a.Name)

	a = r.Asset("windows", "amd64")
	assert.NotNil(t, a, "nil for windows")
	assert.Equal(t, "apex_windows_amd64.exe", a.Name)

	a = r.Asset("sloth", "amd64")
	assert.Nil(t, a)
}

func TestAsset_Download(t *testing.T) {
	t.Parallel()

	p := &Project{
		Command: "apex",
		Owner:   "apex",
		Repo:    "apex",
		Version: "0.13.1",
	}

	releases, err := p.LatestReleases()
	assert.NoError(t, err)
	assert.NotNil(t, releases, "releases")
	r := releases[0]

	a := r.Asset("darwin", "386")
	assert.NotNil(t, a, "nil for darwin")

	path, err := a.Download()
	assert.NoError(t, err, "download")
	assert.NotEmpty(t, path, "path")
}

func TestProject_Install(t *testing.T) {
	t.Parallel()

	p := &Project{
		Command: "apex",
		Owner:   "apex",
		Repo:    "apex",
		Version: "0.13.1",
	}

	releases, err := p.LatestReleases()
	assert.NoError(t, err)
	assert.NotNil(t, releases, "releases")
	r := releases[0]

	a := r.Asset("darwin", "386")
	assert.NotNil(t, a, "nil for darwin")

	path, err := a.Download()
	assert.NoError(t, err, "download")
	assert.NotEmpty(t, path, "path")

	err = p.Install(path)
	assert.NoError(t, err, "replace")
}
