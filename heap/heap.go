package heap

import "container/heap"

type Heap[E any] interface {
	Push(e E)
	Pop() E
	Len() int
}

func NewCmp[S ~[]E, E any](slice S, cmp func(x, y E) int) Heap[E] {
	return newHeap[S](&cmpHeap[S, E]{baseHeap: &baseHeap[S, E]{slice: slice}, cmp: cmp})
}

func NewLess[S ~[]E, E any](slice S, less func(x, y E) bool) Heap[E] {
	return newHeap[S](&lessHeap[S, E]{baseHeap: &baseHeap[S, E]{slice: slice}, less: less})
}

func newHeap[S ~[]E, E any](iface heap.Interface) Heap[E] {
	heap.Init(iface)
	return &heapImpl[E]{Interface: iface}
}

type heapImpl[E any] struct {
	heap.Interface
}

func (h *heapImpl[E]) Push(e E) {
	heap.Push(h.Interface, e)
}

func (h *heapImpl[E]) Pop() E {
	return heap.Pop(h.Interface).(E)
}

type cmpHeap[S ~[]E, E any] struct {
	*baseHeap[S, E]
	cmp func(x, y E) int
}

func (h *cmpHeap[S, E]) Less(i, j int) bool {
	return h.cmp(h.slice[i], h.slice[j]) < 0
}

type lessHeap[S ~[]E, E any] struct {
	*baseHeap[S, E]
	less func(x, y E) bool
}

func (h *lessHeap[S, E]) Less(i, j int) bool {
	return h.less(h.slice[i], h.slice[j])
}

type baseHeap[S ~[]E, E any] struct {
	slice S
}

func (h *baseHeap[S, E]) Len() int {
	return len(h.slice)
}

func (h *baseHeap[S, E]) Swap(i, j int) {
	h.slice[i], h.slice[j] = h.slice[j], h.slice[i]
}

func (h *baseHeap[S, E]) Push(x any) {
	h.slice = append(h.slice, x.(E))
}

func (h *baseHeap[S, E]) Pop() any {
	n := len(h.slice)
	x := h.slice[n-1]
	h.slice = h.slice[:n-1]
	return x
}
