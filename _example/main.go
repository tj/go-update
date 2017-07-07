package main

import (
	"fmt"
	"runtime"

	"github.com/apex/log"
	"github.com/kierdavis/ansi"

	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
)

func init() {
	// log.SetLevel(log.DebugLevel)
}

func main() {
	ansi.HideCursor()
	defer ansi.ShowCursor()

	// update polls(1) from tj/gh-polls on github
	p := &update.Project{
		Command: "polls",
		Owner:   "tj",
		Repo:    "gh-polls",
		Version: "0.0.3",
	}

	// fetch the new releases
	releases, err := p.LatestReleases()
	if err != nil {
		log.Fatalf("error fetching releases: %s", err)
	}

	// no updates
	if len(releases) == 0 {
		log.Info("no updates")
		return
	}

	// latest release
	latest := releases[0]

	// find the tarball for this system
	a := latest.FindTarball(runtime.GOOS, runtime.GOARCH)
	if a == nil {
		log.Info("no binary for your system")
		return
	}

	// whitespace
	fmt.Println()

	// download tarball to a tmp dir
	tarball, err := a.DownloadProxy(progress.Reader)
	if err != nil {
		log.Fatalf("error downloading: %s", err)
	}

	// install it
	if err := p.Install(tarball); err != nil {
		log.Fatalf("error installing: %s", err)
	}

	fmt.Printf("Updated to %s\n", latest.Version)
}
