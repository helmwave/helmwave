package registry

import "golang.org/x/exp/slices"

// IndexOfHost searches registry in slice of registries by host. Returns offset.
func IndexOfHost(a []Config, host string) (int, bool) {
	i := slices.IndexFunc(a, func(r Config) bool {
		return host == r.Host()
	})

	return i, i != -1
}
