package opensearchtools

// PtrTo is a generic constraint that restricts value to be pointers.
// T can be any type.
type PtrTo[T any] interface {
	*T
}
