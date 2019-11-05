package manager

type Set struct {
	s map[string]struct{}
}

func NewSet() *Set {
	return &Set{
		s: map[string]struct{}{},
	}
}

func (s *Set) Add(item string) {
	s.s[item] = struct{}{}
}

func (s *Set) Remove(item string) {
	_, ok := s.s[item]
	if ok {
		delete(s.s, item)
	}
}

func (s *Set) Values() []string {
	res := make([]string, 0)
	for k, _ := range s.s {
		res = append(res, k)
	}
	return res
}
