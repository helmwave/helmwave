package repo

func (rep *Config) In(a []Config) bool {
	for _, r := range a {
		if rep == &r {
			return true
		}
	}

	return false
}

func (rep *Config) InByName(a []Config) bool {
	for _, r := range a {
		if rep.Name == r.Name {
			return true
		}
	}

	return false
}
