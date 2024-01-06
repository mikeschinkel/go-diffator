package diffator

type fixer interface {
	Fixer()
	String() string
}

type Comparator interface {
	Compare() string
}
