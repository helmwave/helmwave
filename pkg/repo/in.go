package repo

// In check that rep in a.
func (rep *config) In(a []Config) bool {
	for i := range a {
		if rep.Name() == a[i].Name() {
			return true
		}
	}

	return false
}

// IndexOf check that rep in a by name.
func (rep *config) IndexOf(a []Config) (int, bool) {
	return IndexOf(a, rep)
}

// IndexOf searches repository in slice of repositories. Returns offset.
func IndexOf(a []Config, rep Config) (i int, found bool) {
	return IndexOfName(a, rep.Name())
}

// IndexOfName searches repository in slice of repositories by name. Returns offset.
func IndexOfName(a []Config, name string) (i int, found bool) {
	for i, r := range a {
		if name == r.Name() {
			return i, true
		}
	}

	return i, false
}
