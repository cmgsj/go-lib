package heap

import "container/heap"

type Heap[T any] interface {
	Push(x T)
	Pop() T
	Len() int
}

func New[S ~[]T, T any](slice S, less func(x, y T) bool) Heap[T] {
	h := &stdHeap[T]{
		slice: slice,
		less:  less,
	}
	heap.Init(h)
	return &heapImpl[T]{h}
}

type heapImpl[T any] struct {
	heap.Interface
}

func (h *heapImpl[T]) Push(x T) { heap.Push(h.Interface, x) }
func (h *heapImpl[T]) Pop() T   { return heap.Pop(h.Interface).(T) }
func (h *heapImpl[T]) Len() int { return h.Interface.Len() }

type stdHeap[T any] struct {
	slice []T
	less  func(x, y T) bool
}

func (h *stdHeap[T]) Len() int           { return len(h.slice) }
func (h *stdHeap[T]) Less(i, j int) bool { return h.less(h.slice[i], h.slice[j]) }
func (h *stdHeap[T]) Swap(i, j int)      { h.slice[i], h.slice[j] = h.slice[j], h.slice[i] }
func (h *stdHeap[T]) Push(x any)         { h.slice = append(h.slice, x.(T)) }
func (h *stdHeap[T]) Pop() any {
	n := len(h.slice)
	x := h.slice[n-1]
	h.slice = h.slice[:n-1]
	return x
}
