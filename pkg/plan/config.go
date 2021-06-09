package plan

type SaveOptions struct {
	file string
	tags []string
	dir  string

	withReleases bool
	withRepos    bool
	withValues   bool
}
