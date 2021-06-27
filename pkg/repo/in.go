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

// InByName check that rep in a by name
func (rep *Config) InByName(a []*Config) bool {
	for _, r := range a {
		if rep.Name == r.Name {
			return true
		}
	}

	return false
}
