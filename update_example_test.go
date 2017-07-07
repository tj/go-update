package update_test

import (
	"fmt"
	"log"
	"runtime"

	"github.com/kierdavis/ansi"

	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
)

func Example() {
	ansi.HideCursor()
	defer ansi.ShowCursor()

	// update apex(1) from apex/apex on github
	p := &update.Project{
		Command: "apex",
		Owner:   "apex",
		Repo:    "apex",
		Version: "0.13.1",
	}

	// fetch the new releases
	releases, err := p.LatestReleases()
	if err != nil {
		log.Fatalf("error fetching releases: %s", err)
	}

	// no updates
	if len(releases) == 0 {
		log.Println("no updates")
		return
	}

	// latest release
	latest := releases[0]

	// find the asset for this system
	a := latest.Asset(runtime.GOOS, runtime.GOARCH)
	if a == nil {
		log.Println("no binary for your system")
		return
	}

	// whitespace
	fmt.Println()

	path, err := a.DownloadProxy(progress.Reader)
	if err != nil {
		log.Fatalf("error downloading: %s", err)
	}

	// replace the previous
	if err := p.Install(path); err != nil {
		log.Fatalf("error installing: %s", err)
	}

	fmt.Printf("Updated to %s\n", latest.Version)
}
