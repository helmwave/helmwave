package release

func NewValuesReference(src, dst string) ValuesReference {
	return ValuesReference{src, dst}
}
