package repo

// In check that rep in a
func (rep *Config) In(a []*Config) bool {
	for _, r := range a {
		if rep == r {
			return true
		}
	}

	return false
}

// IndexOf check that rep in a by name
func (rep *Config) IndexOf(a []*Config) (int, bool) {
	return IndexOf(a, rep)
}

func IndexOf(a []*Config, rep *Config) (i int, found bool) {
	return IndexOfName(a, rep.Name)
}

func IndexOfName(a []*Config, name string) (i int, found bool) {
	for i, r := range a {
		if name == r.Name {
			return i, true
		}
	}

	return i, false
}
