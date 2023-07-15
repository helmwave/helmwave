package repo

import "golang.org/x/exp/slices"

// Equal checks repo configs to have equal names.
func (rep *config) Equal(a Config) bool {
	return rep.Name() == a.Name()
}

// IndexOfName searches repository in slice of repositories by name. Returns offset.
func IndexOfName(a []Config, name string) (int, bool) {
	i := slices.IndexFunc(a, func(r Config) bool {
		return name == r.Name()
	})

	return i, i != -1
}
