package functional

type op int

//goland:noinspection GoUnusedConst
const (
	// opUnknown is used to ensure that the default value is never interpreted as a valid last operation.
	opUnknown = iota

	// opOrderBy indicates that the last operation was an OrderBy[Key][Descending] operation.
	opOrderBy

	// opWhere indicates that the last operation was a Where operation.
	opWhere

	// opSelect indicates that the last operation was a Select operation.
	opSelect
)

type hints[E any] struct {
	lengthEstimate int // used for preallocation, can be inaccurate
	slice          []E
	op             op
}

func newHints[E any]() hints[E] {
	return hints[E]{
		lengthEstimate: -1,
	}
}

func (h hints[E]) clone() hints[E] {
	return hints[E]{
		slice:          h.slice,
		lengthEstimate: h.lengthEstimate,
	}
}

func (h hints[E]) clone2[E2 any]() hints[E2] {
	return hints[E2]{
		lengthEstimate: h.lengthEstimate,
	}
}
