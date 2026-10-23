package functional

func (f *Functioner[E]) Contains(e E, equals func(e1, e2 E) bool) bool {
	return f.Any(func(ee E) bool {
		return equals(e, ee)
	})
}
