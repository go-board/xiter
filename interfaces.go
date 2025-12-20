package xiter

import "iter"

// Extendable describes types that can extend themselves from a sequence.
type Extendable[E any] interface {
	Extend(it iter.Seq[E])
}

// Extendable2 describes types that can extend themselves from a key/value
// sequence.
type Extendable2[K, V any] interface {
	Extend2(it iter.Seq2[K, V])
}
