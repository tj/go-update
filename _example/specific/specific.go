package main

import (
	"fmt"
	"runtime"

	"github.com/apex/log"
	"github.com/kierdavis/ansi"

	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
	"github.com/tj/go-update/stores/github"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	ansi.HideCursor()
	defer ansi.ShowCursor()

	// update polls(1) from tj/gh-polls on github
	m := &update.Manager{
		Command: "up",
		Store: &github.Store{
			Owner:   "apex",
			Repo:    "up",
			Version: "0.4.6",
		},
	}

	// fetch the target release
	release, err := m.GetRelease("0.4.5")
	if err != nil {
		log.Fatalf("error fetching release: %s", err)
	}

	// find the tarball for this system
	a := release.FindTarball(runtime.GOOS, runtime.GOARCH)
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
	if err := m.Install(tarball); err != nil {
		log.Fatalf("error installing: %s", err)
	}

	fmt.Printf("Updated to %s\n", release.Version)
}
