package manager

type Set struct {
	s map[interface{}]struct{}
}

func NewSet() *Set {
	return &Set{
		s: map[interface{}]struct{}{},
	}
}

func (s *Set) Add(item interface{}) {
	s.s[item] = struct{}{}
}

func (s *Set) Remove(item interface{}) {
	_, ok := s.s[item]
	if ok {
		delete(s.s, item)
	}
}

func (s *Set) Values() []interface{} {
	var res []interface{}
	for k, _ := range s.s {
		res = append(res, k)
	}
	return res
}
