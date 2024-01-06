package diffator

type Diffator struct {
	Comparator
}

func NewDiffator(c Comparator) *Diffator {
	d := &Diffator{}
	d.Comparator = c
	return d
}
