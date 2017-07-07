package update_test

import (
	"fmt"
	"log"
	"runtime"

	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
)

func Example() {
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
		log.Println("no updates")
		return
	}

	// latest release
	latest := releases[0]

	// find the asset for this system
	a := latest.FindTarball(runtime.GOOS, runtime.GOARCH)
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
