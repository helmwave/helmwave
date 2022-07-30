package release

func NewValuesReference(src, dst string) ValuesReference {
	return ValuesReference{
		Src: src,
		Dst: dst,
	}
}
