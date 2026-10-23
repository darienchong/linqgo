package functional

import "fmt"

type Pair[K, V any] struct {
	First  K
	Second V
}

func (p *Pair[K, V]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}
