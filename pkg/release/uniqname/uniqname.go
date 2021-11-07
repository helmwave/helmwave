package uniqname

type UniqName string

func Contains(t UniqName, a []UniqName) bool {
	for _, v := range a {
		if v == t {
			return true
		}
	}

	return false
}
