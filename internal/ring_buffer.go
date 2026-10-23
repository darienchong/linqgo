package internal

type RingBuffer[E any] struct {
	elements []E
	cap      int
	start    int
	end      int
}

func NewRingBuffer[E any](cap int) *RingBuffer[E] {
	return &RingBuffer[E]{
		elements: make([]E, 0, cap),
		cap:      cap,
		start:    -1,
	}
}

func (b *RingBuffer[E]) Push(e E) {
	if b.start == -1 {
		b.elements = append(b.elements, e)
		b.start = 0
		return
	}

	if len(b.elements) < b.cap {
		b.elements = append(b.elements, e)
		return
	}

	b.elements[b.start] = e
	b.start = (b.start + 1) % b.cap
}

func (b *RingBuffer[E]) Iterate(fn func(E)) {
	if b.start == -1 {
		return
	}
	for offset := 0; offset < len(b.elements); offset++ {
		fn(b.elements[((b.start + offset) % b.cap)])
	}
}

func (b *RingBuffer[E]) Iterate2(fn func(E) bool) {
	if b.start == -1 {
		return
	}
	for offset := 0; offset < len(b.elements); offset++ {
		if !fn(b.elements[((b.start + offset) % b.cap)]) {
			return
		}
	}
}

func (b *RingBuffer[E]) ReverseIterate2(fn func(E) bool) {
	if b.start == -1 {
		return
	}
	for offset := len(b.elements) - 1; offset >= 0; offset-- {
		if !fn(b.elements[((b.start + offset) % b.cap)]) {
			return
		}
	}
}

func (b *RingBuffer[E]) Collect() []E {
	res := make([]E, 0, len(b.elements))
	for i := 0; i < len(b.elements); i++ {
		res = append(res, b.elements[(b.start+i)%b.cap])
	}
	return res
}

func (b *RingBuffer[E]) Len() int {
	return len(b.elements)
}
