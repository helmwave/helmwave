package repo

func (rep *Config) In(a []Config) bool {
	for _, r := range a {
		if rep == &r {
			return true
		}
	}
	return false
}
