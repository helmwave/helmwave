package repo

import "slices"

// IndexOfName searches repository in slice of repositories by name. Returns offset.
func IndexOfName(a []Config, name string) (int, bool) {
	i := slices.IndexFunc(a, func(r Config) bool {
		return name == r.Name()
	})

	return i, i != -1
}
