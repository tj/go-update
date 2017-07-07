// Package update provides tooling to auto-update binary releases
// from GitHub based on the user's current version and operating system.
package update

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/c4milo/unpackit"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// Proxy is used to proxy a reader, for example
// using https://github.com/cheggaaa/pb to provide
// progress updates.
type Proxy func(int, io.ReadCloser) io.ReadCloser

// NopProxy does nothing.
var NopProxy = func(size int, r io.ReadCloser) io.ReadCloser {
	return r
}

// Project represents the project.
type Project struct {
	Command string // Command is the executable's name.
	Owner   string // Owner is the GitHub owner name.
	Repo    string // Repo is the GitHub repo name.
	Version string // Version is the local version.
}

// Release represents a project release.
type Release struct {
	p           *Project  // Project is the parent project.
	Version     string    // Version is the release version.
	Notes       string    // Notes is the markdown release notes.
	URL         string    // URL is the notes url.
	PublishedAt time.Time // PublishedAt is the publish time.
	Assets      []*Asset  // Assets is the release assets.
}

// Asset represents a project release asset.
type Asset struct {
	Name      string // Name of the asset.
	Size      int    // Size of the asset.
	URL       string // URL of the asset.
	Downloads int    // Downloads count.
}

// InstallTo binary to the given dir.
func (p *Project) InstallTo(path, dir string) error {
	log.Debugf("unpacking %q", path)

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "opening tarball")
	}

	tmpdir, err := unpackit.Unpack(f, "")
	if err != nil {
		f.Close()
		return errors.Wrap(err, "unpacking tarball")
	}

	if err := f.Close(); err != nil {
		return errors.Wrap(err, "closing tarball")
	}

	bin := filepath.Join(tmpdir, p.Command)

	if err := os.Chmod(bin, 0755); err != nil {
		return errors.Wrap(err, "chmod")
	}

	dst := filepath.Join(dir, p.Command)

	log.Debugf("move %q to %q", bin, dst)
	if err := os.Rename(bin, dst); err != nil {
		return errors.Wrap(err, "moving")
	}

	return nil
}

// Install binary to replace the current version.
func (p *Project) Install(path string) error {
	bin, err := exec.LookPath(p.Command)
	if err != nil {
		return errors.Wrapf(err, "looking up path of %q", p.Command)
	}

	dir := filepath.Dir(bin)
	return p.InstallTo(path, dir)
}

// LatestReleases returns releases newer than Version, or nil.
func (p *Project) LatestReleases() (latest []*Release, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gh := github.NewClient(nil)

	releases, _, err := gh.Repositories.ListReleases(ctx, p.Owner, p.Repo, nil)
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		tag := r.GetTagName()

		if tag == p.Version || "v"+p.Version == tag {
			break
		}

		latest = append(latest, toRelease(p, r))
	}

	return
}

// FindTarball returns a tarball matching os and arch, or nil.
func (r *Release) FindTarball(os, arch string) *Asset {
	s := fmt.Sprintf("%s_%s", os, arch)
	for _, a := range r.Assets {
		ext := filepath.Ext(a.Name)
		if strings.Contains(a.Name, s) && ext == ".gz" {
			return a
		}
	}

	return nil
}

// Download the asset to a tmp directory and return its path.
func (a *Asset) Download() (string, error) {
	return a.DownloadProxy(NopProxy)
}

// DownloadProxy the asset to a tmp directory and return its path.
func (a *Asset) DownloadProxy(proxy Proxy) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "update-")
	if err != nil {
		return "", errors.Wrap(err, "creating temp file")
	}

	log.Debugf("fetch %q", a.URL)
	res, err := http.Get(a.URL)
	if err != nil {
		return "", errors.Wrap(err, "fetching asset")
	}

	kind := res.Header.Get("Content-Type")
	size, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	log.Debugf("response %s – %s (%d KiB)", res.Status, kind, size/1024)

	body := proxy(size, res.Body)

	if res.StatusCode >= 400 {
		body.Close()
		return "", errors.Wrap(err, res.Status)
	}

	log.Debugf("copy to %q", f.Name())
	if _, err := io.Copy(f, body); err != nil {
		body.Close()
		return "", errors.Wrap(err, "copying body")
	}

	if err := body.Close(); err != nil {
		return "", errors.Wrap(err, "closing body")
	}

	if err := f.Close(); err != nil {
		return "", errors.Wrap(err, "closing file")
	}

	log.Debugf("copied")
	return f.Name(), nil
}

// toRelease returns a Release.
func toRelease(p *Project, r *github.RepositoryRelease) *Release {
	out := &Release{
		p:           p,
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
