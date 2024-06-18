package netpool

type Router interface {
	Next(current, total uint32) uint32
}

func RoundRobin() Router {
	return &roundRobinRouter{}
}

func Sticky(times uint32) Router {
	return &stickyRouter{times: times}
}

type roundRobinRouter struct{}

func (r *roundRobinRouter) Next(current, total uint32) uint32 {
	return (current + 1) % total
}

type stickyRouter struct {
	count uint32
	times uint32
}

func (s *stickyRouter) Next(current, total uint32) uint32 {
	if s.times == 0 {
		s.times = 1
	}
	if s.count == s.times {
		s.count = 0
		current = (current + 1) % total
	}
	s.count++
	return current
}
