package update

// TODO: tweak... interface should work for
// downgrades to specific versions as well

// Store is the interface used for listing and fetching releases.
type Store interface {
	LatestReleases() ([]*Release, error)
}
