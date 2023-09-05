package registry

// IndexOfHost searches registry in slice of registries by host. Returns offset.
func IndexOfHost(a []Config, host string) (i int, found bool) {
	for i, r := range a {
		if host == r.Host() {
			return i, true
		}
	}

	return i, false
}
